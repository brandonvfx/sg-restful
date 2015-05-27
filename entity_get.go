package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func entityGetHandler(rw http.ResponseWriter, req *http.Request) {
	log.Debug("Calling entityGetHandler")
	vars := mux.Vars(req)
	entity_type := vars["entity_type"]
	log.Debug("Entity Type:", entity_type)

	query := map[string]interface{}{
		"return_fields":         nil,
		"type":                  entity_type,
		"return_paging_info":    true,
		"api_return_image_urls": true,
		"return_only":           "active",
		"paging": map[string]int{
			"current_page":      1,
			"entities_per_page": 1,
		},
		"filters": nil,
	}

	entityIdStr, hasId := vars["id"]
	if hasId {
		entityId, err := strconv.Atoi(entityIdStr)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		query["filters"] = map[string]interface{}{
			"logical_operator": "and",
			"conditions": []map[string]interface{}{
				map[string]interface{}{
					"path":     "id",
					"relation": "is",
					"values":   []int{int(entityId)},
				},
			},
		}
	} else {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(rw, "Id missing")
		return
	}

	req.ParseForm()

	fieldsStr := req.FormValue("fields")
	fields := []string{"id"}
	if fieldsStr != "" {
		fields = strings.Split(fieldsStr, ",")
	}
	query["return_fields"] = fields

	log.Debug(query)

	sg_conn, ok := context.GetOk(req, "sg_conn")
	if !ok {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	sg := sg_conn.(Shotgun)
	sgReq, err := sg.Request("read", query)
	if err != nil {
		log.Error("Request Error: ", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	var readResp readResponse
	respBody, err := ioutil.ReadAll(sgReq.Body)
	if err != nil {
		log.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(respBody, &readResp)
	if err != nil {
		log.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Debug("Response: ", readResp)

	if len(readResp.Results.Entities) == 0 {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	jsonResp, err := json.Marshal(readResp.Results.Entities[0])

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(jsonResp)
}
