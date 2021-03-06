package main

import (
	"net/http"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

var queryRegexp = regexp.MustCompile(`^([\w]+)\((.*)\)`)

type queryParseError struct {
	StatusCode int
	Message    string
}

func (qpe queryParseError) Error() string {
	return qpe.Message
}

func parseQuery(queryStr string) (readFilters, error) {
	manager := GetQPManager()
	keys, parsers := manager.GetActiveParsers()
	for idx, parser := range parsers {
		if parser.CanParseString(queryStr) {
			log.Debugf("parser.CanParseString(%s) true for %s\n", queryStr, keys[idx])
			rf, err := parser.ParseString(queryStr)
			if err != nil {
				log.Warnf("Parser.ParseString(%s) failed with:%s", queryStr, err.Error())
				return rf, queryParseError{
					StatusCode: http.StatusBadRequest,
					Message:    err.Error(),
				}
			}
			return rf, nil
		}
		log.Debugf("parser.CanParseString(%s) false for %s\n", queryStr, keys[idx])

	}

	log.Warnf("parseQuery - Failed to find parser for:%s. Tried %s\n", queryStr, keys)

	return newReadFilters(), queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "No QueryFormats can parse input",
	}
}

//
// func parseQuery(queryStr string) (readFilters, error) {
// 	query := newReadFilters()
//
// 	if strings.HasPrefix(queryStr, "AND") || strings.HasPrefix(queryStr, "and") ||
// 		strings.HasPrefix(queryStr, "OR") || strings.HasPrefix(queryStr, "or") {
// 		// Format: and([key, comparitor, value],[key, comparitor, value] ...)
// 		//          or([key, comparitor, value],[key, comparitor, value] ...)
// 		log.Debugf("Starts with AND/OR")
// 		matches := queryRegexp.FindStringSubmatch(queryStr)
// 		if matches != nil {
// 			if len(matches) != 3 {
// 				return query, queryParseError{
// 					StatusCode: http.StatusBadRequest,
// 					Message:    "Invalid query format",
// 				}
// 			}
//
// 			var filtersStr string
// 			if strings.HasPrefix(matches[2], "[[") {
// 				filtersStr = matches[2]
// 			} else if strings.HasPrefix(matches[2], "[") {
// 				filtersStr = "[" + matches[2] + "]"
// 			} else {
// 				return query, queryParseError{
// 					StatusCode: http.StatusBadRequest,
// 					Message:    "Invalid query filter format",
// 				}
// 			}
//
// 			var filters []interface{}
// 			err := json.Unmarshal([]byte(filtersStr), &filters)
// 			if err != nil {
// 				return query, queryParseError{
// 					StatusCode: http.StatusBadRequest,
// 					Message:    err.Error(),
// 				}
// 			}
//
// 			queryData := map[string]interface{}{
// 				"logical_operator": matches[1],
// 				"conditions":       filters,
// 			}
//
// 			err = mapToReadFilters(queryData, &query)
//
// 			return query, err
//
// 		}
// 		return query, queryParseError{
// 			StatusCode: http.StatusBadRequest,
// 			Message:    "Invalid query format",
// 		}
//
// 	} else if strings.HasPrefix(queryStr, "{") {
// 		// Format: {"logical_operator": "and", "conditions": [[key, comparitor, value], ...]}
// 		log.Debugf("Is JSON hash")
// 		var queryData map[string]interface{}
// 		err := json.Unmarshal([]byte(queryStr), &queryData)
// 		if err != nil {
// 			return query, queryParseError{
// 				StatusCode: http.StatusBadRequest,
// 				Message:    err.Error(),
// 			}
// 		}
//
// 		if _, ok := queryData["logical_operator"]; !ok {
// 			return query, queryParseError{
// 				StatusCode: http.StatusBadRequest,
// 				Message:    "Missing key: 'logical_operator'",
// 			}
// 		}
//
// 		if _, ok := queryData["conditions"]; !ok {
// 			return query, queryParseError{
// 				StatusCode: http.StatusBadRequest,
// 				Message:    "Missing key: 'conditions'",
// 			}
// 		}
// 		err = mapToReadFilters(queryData, &query)
//
// 		return query, err
//
// 	} else if strings.HasPrefix(queryStr, "[[") {
// 		// Format: [[key, comparitor, value], ...]
// 		log.Debugf("Is JSON array of arrays")
// 		var filters []interface{}
// 		err := json.Unmarshal([]byte(queryStr), &filters)
// 		if err != nil {
// 			return query, queryParseError{
// 				StatusCode: http.StatusBadRequest,
// 				Message:    err.Error(),
// 			}
// 		}
// 		queryData := map[string]interface{}{
// 			"logical_operator": "and",
// 			"conditions":       filters,
// 		}
// 		err = mapToReadFilters(queryData, &query)
// 		return query, err
//
// 	}
//
// 	return query, queryParseError{
// 		StatusCode: http.StatusBadRequest,
// 		Message:    "Invalid query format",
// 	}
//
// }
//
func mapToReadFilters(queryMap map[string]interface{}, query *readFilters) error {

	op, ok := queryMap["logical_operator"]
	if !ok {
		return queryParseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Missing key: 'logical_operator'",
		}
	}

	conditions, ok := queryMap["conditions"]
	if !ok {
		return queryParseError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Missing key: 'conditions'",
		}
	}

	query.LogicalOperator = strings.ToLower(op.(string))

	log.Debugf("Conditions: %v", conditions)
	for _, condition := range conditions.([]interface{}) {
		filter := condition.([]interface{})
		log.Debugf("Filter: %v", filter)
		cond := arrayToQueryCondition(filter)
		log.Debugf("Condition: %v", cond)
		query.AddCondition(cond)
	}
	return nil
}

func arrayToQueryCondition(filter []interface{}) queryCondition {
	return newQueryCondition(
		filter[0].(string), // name
		filter[1].(string), // relation
		filter[2],          // values
	)
}
