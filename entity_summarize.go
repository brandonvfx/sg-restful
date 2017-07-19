package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Response Structs
type summaryResponse struct {
	Summaries map[string]interface{} `json:"summaries"`
	Groups    []groupResponse        `json:"groups"`
}

type groupResponse struct {
	GroupValue string                 `json:"group_value"`
	GroupName  string                 `json:"group_name"`
	Summaries  map[string]interface{} `json:"summaries"`
	Groups     []groupResponse        `json:"groups,omitempty"`
}

type summarizeResponse struct {
	Results   summaryResponse `json:"results"`
	Exception bool            `json:"exception,omitempty"`
	Message   string          `json:"message,omitempty"`
	ErrorCode int             `json:"error_code,omitempty"`
}

// Query Structs
type summarizeQuery struct {
	EntityType string      `json:"type"`
	Filters    readFilters `json:"filters"`
	Grouping   []grouping  `json:"grouping,omitempty"`
	Summaries  []summary   `json:"summaries,omitempty"`
}

func newSummarizeQuery(entityType string) summarizeQuery {
	return summarizeQuery{
		EntityType: entityType,
		Filters:    newReadFilters(),
		Summaries:  make([]summary, 0),
	}
}

type grouping struct {
	Direction string `json:"direction"`
	Field     string `json:"field"`
	Type      string `json:"type"`
}

type summary struct {
	Field string `json:"field"`
	Type  string `json:"type"`
}

// Handlers

func entitySummarizeHandler(config clientConfig) func(rw http.ResponseWriter, req *http.Request) {
	return func(rw http.ResponseWriter, req *http.Request) {
		log.Debug("Calling entitySummarizeHandler")
		vars := mux.Vars(req)
		entityType, ok := vars["entity_type"]
		if !ok {
			log.Errorf("Missing Entity Type")
			rw.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Debugf("Entity: %s", entityType)

		query := newSummarizeQuery(entityType)

		req.ParseForm()

		// Since there would be any number of "fields" on an entity
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
			case "summaries":
				var summaries []summary
				err := json.Unmarshal(bytes.NewBufferString(value).Bytes(), &summaries)
				if err != nil {
					log.Error(err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				log.Debugf("Summary: %v", summaries)
				query.Summaries = summaries
			case "grouping":
				var groups []grouping
				err := json.Unmarshal(bytes.NewBufferString(value).Bytes(), &groups)
				if err != nil {
					log.Error(err)
					rw.WriteHeader(http.StatusInternalServerError)
					return
				}
				log.Debugf("Groups: %v", groups)
				query.Grouping = groups
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

		sgReq, err := sg.Request("summarize", query)
		if err != nil {
			log.Error("Request Error: ", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		var summarizeResp summarizeResponse
		respBody, err := ioutil.ReadAll(sgReq.Body)
		if err != nil {
			log.Error(err)
			rw.WriteHeader(http.StatusBadGateway)
			return
		}

		err = json.Unmarshal(respBody, &summarizeResp)
		if err != nil {
			log.Error(err)
			rw.WriteHeader(http.StatusBadGateway)
			return
		}

		log.Debugf("Response: %v", summarizeResp)

		if summarizeResp.Exception {
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(summarizeResp.Message))
			return
		}

		// if len(summarizeResp.Results.Groups) == 0 {
		// 	rw.WriteHeader(http.StatusNoContent)
		// 	return
		// }

		var jsonResp []byte
		jsonResp, err = json.Marshal(summarizeResp.Results)

		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(jsonResp)
	}
}
