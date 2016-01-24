package main

/*
Implementation of QueryParserI interface
*/

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// Format1 satisfies the QueryParserI interface for format1
// from the sg-query spec
type Format1 struct{}

var format1QueryRegexp = regexp.MustCompile(`^([\w]+)\((.*)\)`)

// CanParseString returns a boolean indicating whether or not the method can parse the supplied string. It is not a guarantee that parsing will be successful.
func (f *Format1) CanParseString(queryStr string) bool {
	// gotta start with and/or
	if !strings.HasPrefix(queryStr, "AND") &&
		!strings.HasPrefix(queryStr, "and") &&
		!strings.HasPrefix(queryStr, "OR") &&
		!strings.HasPrefix(queryStr, "or") {
		return false
	}

	matches := format1QueryRegexp.FindStringSubmatch(queryStr)
	if matches != nil {
		if len(matches) != 3 {
			return false
		}

		if strings.HasPrefix(matches[2], "[[") || strings.HasPrefix(matches[2], "[") {
			return true
		}
	}
	return false
}

// ParseString parses the supplied string and returns readFilters and an error
func (f *Format1) ParseString(queryStr string) (readFilters, error) {
	query := newReadFilters()

	if !strings.HasPrefix(queryStr, "AND") &&
		!strings.HasPrefix(queryStr, "and") &&
		!strings.HasPrefix(queryStr, "OR") &&
		!strings.HasPrefix(queryStr, "or") {
		return query, queryParseError{
			StatusCode: http.StatusBadRequest,
			Message:    "Invalid query format",
		}
	}

	matches := queryRegexp.FindStringSubmatch(queryStr)
	if matches != nil {
		if len(matches) != 3 {
			return query, queryParseError{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid query format",
			}
		}

		var filtersStr string
		if strings.HasPrefix(matches[2], "[[") {
			filtersStr = matches[2]
		} else if strings.HasPrefix(matches[2], "[") {
			filtersStr = "[" + matches[2] + "]"
		} else {
			return query, queryParseError{
				StatusCode: http.StatusBadRequest,
				Message:    "Invalid query filter format",
			}
		}

		var filters []interface{}
		err := json.Unmarshal([]byte(filtersStr), &filters)
		if err != nil {
			return query, queryParseError{
				StatusCode: http.StatusBadRequest,
				Message:    err.Error(),
			}
		}

		queryData := map[string]interface{}{
			"logical_operator": matches[1],
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

// Register the format with the manager
func init() {
	manager := GetQPManager()
	log.Info("manager fetched", manager)
	manager.AddParser("format1", &Format1{})
}
