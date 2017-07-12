package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type updateResponse struct {
	Results   map[string]interface{} `json:"results"`
	Exception bool                   `json:"exception,omitempty"`
	Message   string                 `json:"message,omitempty"`
	ErrorCode int                    `json:"error_code,omitempty"`
}

func entityUpdateHandler(config clientConfig) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		entityType, ok := vars["entity_type"]
		if !ok {
			log.Errorf("Missing Entity Type")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		var entityID int
		var err error
		entityIDStr, ok := vars["id"]
		if ok {
			entityID, err = strconv.Atoi(entityIDStr)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Debugf("Entity: %s - %d", entityType, entityID)

		var patchData map[string]interface{}
		patchBody, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Error(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		err = json.Unmarshal(patchBody, &patchData)
		if err != nil {
			log.Error(err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Info("Patch Data:", patchData)

		fields := make([]map[string]interface{}, len(patchData))
		i := 0
		for key, value := range patchData {
			field := make(map[string]interface{})
			field["field_name"] = key
			field["value"] = value
			fields[i] = field
			i++
		}

		query := map[string]interface{}{
			"type":   entityType,
			"id":     entityID,
			"fields": fields,
		}

		ctx := req.Context()
		sgConn := ctx.Value("sgConn")
		if sgConn == nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		sg := sgConn.(Shotgun)
		sgReq, err := sg.Request("update", query)
		if err != nil {
			log.Error("Request Error: ", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		var updateResp updateResponse
		respBody, err := ioutil.ReadAll(sgReq.Body)
		if err != nil {
			log.Error(err)
			rw.WriteHeader(http.StatusBadGateway)
			return
		}

		err = json.Unmarshal(respBody, &updateResp)
		if err != nil {
			log.Error(err)
			rw.WriteHeader(http.StatusBadGateway)
			return
		}

		log.Debug("Response: ", updateResp)

		if updateResp.Exception {
			if strings.Contains(updateResp.Message, "unique") {
				rw.WriteHeader(http.StatusConflict)
			} else if strings.Contains(updateResp.Message, "Permission") {
				rw.WriteHeader(http.StatusForbidden)
			} else if strings.Contains(updateResp.Message, "does not exist") {
				rw.WriteHeader(http.StatusNotFound)
			} else {
				rw.WriteHeader(http.StatusBadRequest)
			}
			rw.Write(bytes.NewBufferString(updateResp.Message).Bytes())
			return
		}

		jsonResp, err := json.Marshal(updateResp.Results)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(jsonResp)
	}
}
