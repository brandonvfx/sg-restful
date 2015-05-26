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

type UpdateResponse struct {
	Results   map[string]interface{} `json:"results"`
	Exception bool                   `json:"exception",omitempty`
	Message   string                 `json:"message",omitempty`
	ErrorCode int                    `json:"error_code",omitempty`
}

func entityUpdateHandler(rw http.ResponseWriter, req *http.Request) {
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

	var patchData map[string]interface{}
	patchBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error(err)
		return
	}

	err = json.Unmarshal(patchBody, &patchData)
	if err != nil {
		log.Error(err)
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
		"type":   entity_type,
		"id":     entityId,
		"fields": fields,
	}

	sg_conn, ok := context.GetOk(req, "sg_conn")
	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	sg := sg_conn.(Shotgun)
	sgReq, err := sg.Request("update", query)
	if err != nil {
		log.Error("Request Error: ", err)
		return
	}

	var updateResp UpdateResponse
	respBody, err := ioutil.ReadAll(sgReq.Body)
	if err != nil {
		log.Error(err)
		return
	}

	err = json.Unmarshal(respBody, &updateResp)
	if err != nil {
		log.Error(err)
		return
	}

	log.Debug("Response: ", updateResp)

	if updateResp.Exception {
		if strings.Contains(updateResp.Message, "unique") {
			rw.WriteHeader(http.StatusConflict)
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
