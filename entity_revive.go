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

type reviveResponse struct {
	Results   bool   `json:"results"`
	Exception bool   `json:"exception,omitempty"`
	Message   string `json:"message,omitempty"`
	ErrorCode int    `json:"error_code,omitempty"`
}

func entityReviveHandler(config clientConfig) func(rw http.ResponseWriter, req *http.Request) {
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
			log.Debugf("Id: %v  Error: %s", entityID, err)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
		} else {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Debugf("Entity: %s - %d", entityType, entityID)

		query := map[string]interface{}{
			"type": entityType,
			"id":   entityID,
		}

		ctx := req.Context()
		sgConn := ctx.Value("sgConn")
		if sgConn == nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		sg := sgConn.(Shotgun)
		sgReq, err := sg.Request("revive", query)
		if err != nil {
			log.Errorf("Request Error: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		var reviveResp reviveResponse
		respBody, err := ioutil.ReadAll(sgReq.Body)
		if err != nil {
			log.Error(err)
			rw.WriteHeader(http.StatusBadGateway)
			return
		}

		log.Debugf("Json Response: %s", respBody)

		err = json.Unmarshal(respBody, &reviveResp)
		if err != nil {
			log.Error(err)
			rw.WriteHeader(http.StatusBadGateway)
			return
		}

		log.Debugf("Response: %v", reviveResp)

		if reviveResp.Exception {
			if strings.Contains(reviveResp.Message, "Permission") {
				rw.WriteHeader(http.StatusForbidden)
			} else if strings.Contains(reviveResp.Message, "does not exist") {
				rw.WriteHeader(http.StatusNotFound)
			} else {
				rw.WriteHeader(http.StatusBadRequest)
			}
			rw.Write(bytes.NewBufferString(reviveResp.Message).Bytes())
			return
		}

		// I'm not sure this can even happen
		if !reviveResp.Results {
			rw.WriteHeader(http.StatusNotFound)
			return
		}

		rw.WriteHeader(http.StatusOK)
	}
}
