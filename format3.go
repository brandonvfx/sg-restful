package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Format3 satisfies the QueryParserI interface for format1
// from the sg-query spec
type Format3 struct{}

// CanParseString returns a boolean indicating whether or not the method can parse the supplied string. It is not a guarantee that parsing will be successful.
func (f *Format3) CanParseString(queryStr string) bool {

	if !strings.HasPrefix(queryStr, "[[") {
		return false
	}

	return true
}

// ParseString parses the supplied string and returns readFilters and an error
func (f *Format3) ParseString(queryStr string) (readFilters, error) {
	query := newReadFilters()
	if strings.HasPrefix(queryStr, "[[") {
		// Format: [[key, comparitor, value], ...]
		var filters []interface{}
		err := json.Unmarshal([]byte(queryStr), &filters)
		if err != nil {
			return query, queryParseError{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			}
		}
		queryData := map[string]interface{}{
			"logical_operator": "and",
			"conditions":       filters,
		}
		err = mapToReadFilters(queryData, &query)
		return query, err

	}

	return query, queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid query format",
	}
}

// init Register the format with the manager
func init() {
	manager := GetQPManager()
	manager.AddParser("format3", &Format3{})
}
