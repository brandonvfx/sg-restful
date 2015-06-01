package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type DeleteResponse struct {
	Results   bool   `json:"results"`
	Exception bool   `json:"exception",omitempty`
	Message   string `json:"message",omitempty`
	ErrorCode int    `json:"error_code",omitempty`
}

func entityDeleteHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	entity_type := vars["entity_type"]
	log.Debugf("Entity Type: %s", entity_type)

	var entityId int
	var err error
	entityIdStr, hasId := vars["id"]
	if hasId {
		entityId, err = strconv.Atoi(entityIdStr)
		log.Debugf("Id: %v  Error: %s", entityId, err)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	query := map[string]interface{}{
		"type": entity_type,
		"id":   entityId,
	}

	sg_conn, ok := context.GetOk(req, "sg_conn")
	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	sg := sg_conn.(Shotgun)
	sgReq, err := sg.Request("delete", query)
	if err != nil {
		log.Error("Request Error: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	var deleteResp DeleteResponse
	respBody, err := ioutil.ReadAll(sgReq.Body)
	if err != nil {
		log.Error(err)
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	log.Debugf("Json Response: %s", respBody)

	err = json.Unmarshal(respBody, &deleteResp)
	if err != nil {
		log.Error(err)
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	log.Debugf("Response: %v", deleteResp)

	if deleteResp.Exception {
		if strings.Contains(deleteResp.Message, "Permission") {
			rw.WriteHeader(http.StatusForbidden)
		} else if strings.Contains(deleteResp.Message, "does not exist") {
			rw.WriteHeader(http.StatusNotFound)
		} else {
			rw.WriteHeader(http.StatusBadRequest)
		}
		rw.Write(bytes.NewBufferString(deleteResp.Message).Bytes())
		return
	}

	// I'm not sure this can even happen
	if !deleteResp.Results {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
