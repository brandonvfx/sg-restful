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

func entityGetHandler(config clientConfig) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		log.Debug("Calling entityGetHandler")
		vars := mux.Vars(req)
		entityType, ok := vars["entity_type"]
		if !ok {
			log.Errorf("Missing Entity Type")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		query := map[string]interface{}{
			"return_fields":         nil,
			"type":                  entityType,
			"return_paging_info":    true,
			"api_return_image_urls": true,
			"return_only":           "active",
			"paging": map[string]int{
				"current_page":      1,
				"entities_per_page": 1,
			},
			"filters": nil,
		}

		entityIDStr, ok := vars["id"]
		if ok {
			entityID, err := strconv.Atoi(entityIDStr)
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
						"values":   []int{int(entityID)},
					},
				},
			}
			log.Debugf("Entity: %s - %d", entityType, entityID)
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

		sgConn, ok := context.GetOk(req, "sgConn")
		if !ok {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		sg := sgConn.(Shotgun)
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
			rw.WriteHeader(http.StatusBadGateway)
			return
		}
		err = json.Unmarshal(respBody, &readResp)
		if err != nil {
			log.Error(err)
			rw.WriteHeader(http.StatusBadGateway)
			return
		}

		log.Debugf("Response: %v", readResp)

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
}
