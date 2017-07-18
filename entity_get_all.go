package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Response Structs
type entityResponse struct {
	Entities   []map[string]interface{} `json:"entities"`
	PagingInfo map[string]int           `json:"paging_info"`
}

type readResponse struct {
	Results   entityResponse `json:"results"`
	Exception bool           `json:"exception,omitempty"`
	Message   string         `json:"message,omitempty"`
	ErrorCode int            `json:"error_code,omitempty"`
}

// Query Structs
type readQuery struct {
	ReturnFields       []string       `json:"return_fields"`
	Type               string         `json:"type"`
	ReturnPagingInfo   bool           `json:"return_paging_info"`
	APIReturnImageUrls bool           `json:"api_return_image_urls"`
	ReturnOnly         string         `json:"return_only"`
	Paging             map[string]int `json:"paging"`
	Filters            readFilters    `json:"filters"`
}

func newReadQuery(entityType string) readQuery {
	return readQuery{
		ReturnFields:       []string{"id"},
		Type:               entityType,
		ReturnPagingInfo:   true,
		APIReturnImageUrls: true,
		ReturnOnly:         "active",
		Paging: map[string]int{
			"current_page":      1,
			"entities_per_page": 500,
		},
		Filters: newReadFilters(),
	}
}

type readFilters struct {
	LogicalOperator string           `json:"logical_operator"`
	Conditions      []queryCondition `json:"conditions"`
}

func newReadFilters() readFilters {
	return readFilters{
		LogicalOperator: "and",
		Conditions:      make([]queryCondition, 0),
	}
}

func (rf *readFilters) AddCondition(cond queryCondition) {
	rf.Conditions = append(rf.Conditions, cond)
}

type queryCondition struct {
	Path     string        `json:"path"`
	Relation string        `json:"relation"`
	Values   []interface{} `json:"values"`
}

func newQueryCondition(path string, relation string, values interface{}) queryCondition {
	cond := queryCondition{
		Path:     path,
		Relation: relation,
	}
	typeof := reflect.TypeOf(values)
	if typeof != nil {
		switch typeof.Kind() {
		case reflect.Slice:
			cond.Values = values.([]interface{})
		default:
			cond.Values = []interface{}{values}
		}
	} else {
		cond.Values = []interface{}{values}
	}
	return cond
}

// Handlers

func entityGetAllHandler(config clientConfig) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		log.Debug("Calling entityGetAllHandler")
		vars := mux.Vars(req)
		entityType, ok := vars["entity_type"]
		if !ok {
			log.Errorf("Missing Entity Type")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Debugf("Entity: %s", entityType)

		query := newReadQuery(entityType)

		req.ParseForm()

		// Since there woulc be any number of "fields" on an entity
		// and we want to allow filtering on thoses via the query string.
		// We have to loop over all query string KVs and pull out the reserved ones
		// and add all others to the filters.
		// NOTE: right now we only support simple filtering 'name=foo' becomes ['name', 'is', 'foo']
		//       I want to add better filtering like ^ for startswith, $ for endswith, % for contains.
		//       For more advanced searching a new endpoint for search will be added.
		for k := range req.Form {
			value := req.FormValue(k)
			log.Debugf("Field: '%v' Value: '%v'", k, value)

			switch strings.ToLower(k) {
			case "page":
				if value != "" {
					page, err := strconv.Atoi(value)
					log.Infof("Could not convert page '%v' to int", value)
					if err != nil {
						log.Errorf("Could not convert page '%v' to int", value)
						rw.WriteHeader(http.StatusBadRequest)
						fmt.Fprintf(rw, "Could not convert page '%v' to int", value)
						return
					}
					query.Paging["current_page"] = page
				}
			case "limit":
				if value != "" {
					limit, err := strconv.Atoi(value)
					log.Infof("Could not convert limit '%v' to int", value)
					if err != nil {
						log.Errorf("Could not convert limit '%v' to int", value)
						rw.WriteHeader(http.StatusBadRequest)
						fmt.Fprintf(rw, "Could not convert limit '%v' to int", value)
						return
					}
					query.Paging["entities_per_page"] = limit

				}
			case "fields":
				fields := []string{"id"}
				if value != "" {
					fields = strings.Split(value, ",")
					query.ReturnFields = fields
				}
			case "q":
				// var queryData [][]interface{}
				queryFilters, err := parseQuery(value)
				if err != nil {
					qpeError := err.(queryParseError)
					log.Error("Request Error: ", qpeError)
					rw.WriteHeader(qpeError.StatusCode)
					return
				}
				query.Filters = queryFilters

				log.Debugf("Query: %s", StructToString(query))
				jsonQuery, err := json.Marshal(query)
				if err != nil {
					log.Error(err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				log.Debugf("query json: %s", jsonQuery)
			}

		}

		log.Debugf("Query: %v", StructToString(query))

		ctx := req.Context()
		sgConn := ctx.Value("sgConn")
		if sgConn == nil {
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

		if readResp.Exception {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(readResp.Message))
			return
		}

		if len(readResp.Results.Entities) == 0 {
			rw.WriteHeader(http.StatusNoContent)
			return
		}

		var jsonResp []byte
		jsonResp, err = json.Marshal(readResp.Results.Entities)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(jsonResp)
	}
}
