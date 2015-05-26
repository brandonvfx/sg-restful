package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type DeleteResponse struct {
	Results bool `json:"results"`
	// Exception bool                   `json:"exception",omitempty`
	// Message   string                 `json:"message",omitempty`
	// ErrorCode int                    `json:"error_code",omitempty`
}

func entityDeleteHandler(rw http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	entity_type := vars["entity_type"]
	log.Debug("Entity Type:", entity_type)

	var entityId int
	var err error
	entityIdStr, hasId := vars["id"]
	if hasId {
		entityId, err = strconv.Atoi(entityIdStr)
		log.Debug(entityId)
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
		return
	}

	var deleteResp DeleteResponse
	respBody, err := ioutil.ReadAll(sgReq.Body)
	if err != nil {
		log.Error(err)
		return
	}

	err = json.Unmarshal(respBody, &deleteResp)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debug("Response: ", deleteResp)

	if !deleteResp.Results {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
