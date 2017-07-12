package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type createResponse struct {
	Results   map[string]interface{} `json:"results"`
	Exception bool                   `json:"exception,omitempty"`
	Message   string                 `json:"message,omitempty"`
	ErrorCode int                    `json:"error_code,omitempty"`
}

func entityCreateHandler(config clientConfig) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		entityType, ok := vars["entity_type"]
		if !ok {
			log.Errorf("Missing Entity Type")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Debugf("Entity: %s", entityType)

		var postData map[string]interface{}
		postBody, err := ioutil.ReadAll(req.Body)
		log.Debugf("Post Body: %s", postBody)
		if err != nil {
			log.Errorf("Bad Request Body: %v", err)
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(rw, err)
			return
		}

		err = json.Unmarshal(postBody, &postData)
		if err != nil {
			log.Errorf("Bad Json: %v", err)
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(rw, err)
			return
		}
		log.Debugf("Post Data: %v", postData)

		fields := make([]map[string]interface{}, len(postData))
		i := 0
		for key, value := range postData {
			field := make(map[string]interface{})
			field["field_name"] = key
			field["value"] = value
			fields[i] = field
			i++
		}

		query := map[string]interface{}{
			"return_fields": []string{"id"},
			"type":          entityType,
			"fields":        fields,
		}

		ctx := req.Context()
		sgConn := ctx.Value("sgConn")
		if sgConn == nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		sg := sgConn.(Shotgun)
		sgReq, err := sg.Request("create", query)
		if err != nil {
			log.Errorf("Request Error: %s", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		var createResp createResponse
		respBody, err := ioutil.ReadAll(sgReq.Body)
		if err != nil {
			log.Error(err)
			rw.WriteHeader(http.StatusBadGateway)
			return
		}
		log.Debugf("Json Response: %s", respBody)

		err = json.Unmarshal(respBody, &createResp)
		if err != nil {
			log.Error(err)
			rw.WriteHeader(http.StatusBadGateway)
			return
		}
		log.Debugf("Response: %s", createResp)

		if createResp.Exception {
			if strings.Contains(createResp.Message, "unique") {
				rw.WriteHeader(http.StatusConflict)
			} else if strings.Contains(createResp.Message, "Permission") {
				rw.WriteHeader(http.StatusForbidden)
			} else {
				rw.WriteHeader(http.StatusBadRequest)
			}
			rw.Write(bytes.NewBufferString(createResp.Message).Bytes())
			return
		}

		jsonResp, err := json.Marshal(createResp.Results)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusCreated)
		rw.Write(jsonResp)
	}
}
