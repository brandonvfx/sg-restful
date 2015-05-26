package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type CreateResponse struct {
	Results   map[string]interface{} `json:"results"`
	Exception bool                   `json:"exception",omitempty`
	Message   string                 `json:"message",omitempty`
	ErrorCode int                    `json:"error_code",omitempty`
}

func entityCreateHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	entity_type := vars["entity_type"]
	log.Debug("Entity Type:", entity_type)

	var postData map[string]interface{}
	postBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error(err)
		return
	}

	err = json.Unmarshal(postBody, &postData)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("Post Data:", postData)

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
		"type":          entity_type,
		"fields":        fields,
	}

	sg_conn, ok := context.GetOk(req, "sg_conn")
	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	sg := sg_conn.(Shotgun)
	sgReq, err := sg.Request("create", query)
	if err != nil {
		log.Error("Request Error: ", err)
		return
	}

	var createResp CreateResponse
	respBody, err := ioutil.ReadAll(sgReq.Body)
	if err != nil {
		log.Error(err)
		return
	}
	err = json.Unmarshal(respBody, &createResp)
	if err != nil {
		log.Error(err)
		return
	}
	log.Debug("Response: ", createResp)

	if createResp.Exception {
		if strings.Contains(createResp.Message, "unique") {
			rw.WriteHeader(http.StatusConflict)
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
