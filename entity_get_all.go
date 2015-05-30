package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// Response Structs
type entityResponse struct {
	Entities   []map[string]interface{} `json:"entities"`
	PagingInfo map[string]int           `json:"paging_info"`
}

type readResponse struct {
	Results   entityResponse `json:"results"`
	Exception bool           `json:"exception",omitempty`
	Message   string         `json:"message",omitempty`
	ErrorCode int            `json:"error_code",omitempty`
}

// Query Structs
type readQuery struct {
	ReturnFields       []string       `json:"return_fields"`
	Type               string         `json:"type"`
	ReturnPagingInfo   bool           `json:"return_paging_info"`
	ApiReturnImageUrls bool           `json:"api_return_image_urls"`
	ReturnOnly         string         `json:"return_only"`
	Paging             map[string]int `json:"paging"`
	Filters            readFilters    `json:"filters"`
}

func newReadQuery(entity_type string) readQuery {
	return readQuery{
		ReturnFields:       []string{"id"},
		Type:               entity_type,
		ReturnPagingInfo:   true,
		ApiReturnImageUrls: true,
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

// Handlers

func entityGetAllHandler(rw http.ResponseWriter, req *http.Request) {
	log.Debug("Calling entityGetAllHandler")
	vars := mux.Vars(req)
	entity_type := vars["entity_type"]
	log.Debug("Entity Type:", entity_type)

	query := newReadQuery(entity_type)

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
				if err != nil {
					log.Errorf("Could not convert page '%v' to int", value)
				} else {
					query.Paging["current_page"] = page
				}
			}
		case "limit":
			if value != "" {
				limit, err := strconv.Atoi(value)
				if err != nil {
					log.Errorf("Could not convert limit '%v' to int", value)
				} else {
					query.Paging["entities_per_page"] = limit
				}
			}
		case "fields":
			fields := []string{"id"}
			if value != "" {
				fields = strings.Split(value, ",")
				query.ReturnFields = fields
			}
		case "q":
			var queryData [][]interface{}

			err := json.Unmarshal([]byte(value), &queryData)
			if err != nil {
				log.Error(err)
				rw.WriteHeader(http.StatusBadRequest)
				return
			}

			log.Debugf("query: %v", queryData)
			jsonQuery, err := json.Marshal(queryData)
			if err != nil {
				log.Error(err)
				rw.WriteHeader(http.StatusBadRequest)
				return
			}
			log.Debugf("query json: %s", jsonQuery)

			for _, filter := range queryData {
				cond := queryCondition{
					Path:     filter[0].(string),
					Relation: filter[1].(string),
					Values:   []interface{}{filter[2]},
				}
				query.Filters.AddCondition(cond)
			}

		default:
			log.Infof("Default: %v, '%v'", k, value)
			var values []interface{}
			qsValues := strings.Split(value, ",")
			log.Infof("Values: %v", qsValues)

			for v := range qsValues {
				valueInt, err := strconv.Atoi(qsValues[v])
				if err == nil {
					log.Infof("Is Int: %v", valueInt)
					values = append(values, valueInt)
					continue
				}
				values = append(values, qsValues[v])
			}

			relation := "is"
			if len(values) > 1 {
				relation = "in"
			}

			cond := queryCondition{
				Path:     k,
				Relation: relation,
				Values:   values,
			}

			query.Filters.AddCondition(cond)
		}

	}

	log.Debugf("Query: %v", query)

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
