package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Format2 satisfies the QueryParserI interface for format1
// from the sg-query spec
type Format2 struct{}

// CanParseString returns a boolean indicating whether or not the method can parse the supplied string. It is not a guarantee that parsing will be successful.
func (f *Format2) CanParseString(queryStr string) bool {
	if !strings.HasPrefix(queryStr, "{") {
		return false
	}
	// Format: {"logical_operator": "and", "conditions": [[key, comparitor, value], ...]}
	var queryData map[string]interface{}
	err := json.Unmarshal([]byte(queryStr), &queryData)
	if err != nil {
		return false
	}

	if _, ok := queryData["logical_operator"]; !ok {
		return false
	}

	if _, ok := queryData["conditions"]; !ok {
		return false
	}
	return true
}

// ParseString parses the supplied string and returns readFilters and an error
func (f *Format2) ParseString(queryStr string) (readFilters, error) {
	query := newReadFilters()
	if strings.HasPrefix(queryStr, "{") {
		// Format: {"logical_operator": "and", "conditions": [[key, comparitor, value], ...]}
		var queryData map[string]interface{}
		err := json.Unmarshal([]byte(queryStr), &queryData)
		if err != nil {
			return query, queryParseError{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			}
		}

		if _, ok := queryData["logical_operator"]; !ok {
			return query, queryParseError{
				StatusCode: http.StatusBadRequest,
				Message:    "Missing key: 'logical_operator'",
			}
		}

		if _, ok := queryData["conditions"]; !ok {
			return query, queryParseError{
				StatusCode: http.StatusBadRequest,
				Message:    "Missing key: 'conditions'",
			}
		}
		err = mapToReadFilters(queryData, &query)

		return query, err
	}
	return query, queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "Missing Prefix: '{'",
	}
}

// Register the format with the manager
func init() {
	manager := GetQPManager()
	manager.AddParser("format2", &Format2{})
}
