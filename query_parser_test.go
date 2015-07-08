package main

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryParserError(t *testing.T) {
	q := queryParseError{
		StatusCode: 400,
		Message:    "Test Message",
	}

	assert.Equal(t, 400, q.StatusCode, "StatusCode Should be 400")
	assert.Equal(t, "Test Message", q.Error(), "Should be 'Test Message'")

}

func TestParseQueryAndOrStatementRegExpError(t *testing.T) {
	_, err := parseQuery("()")

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid query format",
	}
	assert.Equal(t, expectedError, err, "Should be formating error")
}

func TestParseQueryAndOrStatementFilterFormatError(t *testing.T) {
	_, err := parseQuery("and(foo, bar)")

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid query filter format",
	}
	assert.Equal(t, expectedError, err, "Should be formating error")
}

func TestParseQueryAndOrStatementJsonError(t *testing.T) {
	_, err := parseQuery("and([foo, bar,])")

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid character 'o' in literal false (expecting 'a')",
	}
	assert.Equal(t, expectedError, err, "Should be formating error")
}

func TestParseQueryAndOrStatementNotAQuery(t *testing.T) {
	_, err := parseQuery("andwell_this_is_bad")

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "Invalid query format",
	}
	assert.Equal(t, expectedError, err, "Should be formating error")
}

func TestParseQueryAndOrStatementSliceOfFilters(t *testing.T) {
	rf, err := parseQuery(`and([["name", "is", "blorg"]])`)

	rfExpected := newReadFilters()
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	assert.Equal(t, nil, err, "Error should be nil")
	assert.Equal(t, rfExpected, rf, "Should be the same")
}

func TestParseQueryAndOrStatementBasicAnd(t *testing.T) {
	rf, err := parseQuery(`and(["name", "is", "blorg"])`)

	rfExpected := newReadFilters()
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	assert.Equal(t, nil, err, "Error should be nil")
	assert.Equal(t, rfExpected, rf, "Should be the same")
}

func TestParseQueryAndOrStatementBasicAndUpper(t *testing.T) {
	rf, err := parseQuery(`AND(["name", "is", "blorg"])`)

	rfExpected := newReadFilters()
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	assert.Equal(t, nil, err, "Error should be nil")
	assert.Equal(t, rfExpected, rf, "Should be the same")
}

func TestParseQueryAndOrStatementBasicOr(t *testing.T) {
	rf, err := parseQuery(`or(["name", "is", "blorg"])`)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "or"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	assert.Equal(t, nil, err, "Error should be nil")
	assert.Equal(t, rfExpected, rf, "Should be the same")
}

func TestParseQueryAndOrStatementBasicOrUpper(t *testing.T) {
	rf, err := parseQuery(`OR(["name", "is", "blorg"])`)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "or"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	assert.Equal(t, nil, err, "Error should be nil")
	assert.Equal(t, rfExpected, rf, "Should be the same")
}

func TestParseQueryHashJsonError(t *testing.T) {
	_, err := parseQuery("{foo}")

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid character 'f' looking for beginning of object key string",
	}
	assert.Equal(t, expectedError, err, "Should a json error")
}

func TestParseQueryHashConditionsOnly(t *testing.T) {
	_, err := parseQuery(`{"conditions": [["name", "is", "blorg"]]}`)

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "Missing key: 'logical_operator'",
	}
	assert.Equal(t, expectedError, err, "Should a json error")
}

func TestParseQueryHashLogicalOpOnly(t *testing.T) {
	_, err := parseQuery(`{"logical_operator": "and"}`)

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "Missing key: 'conditions'",
	}
	assert.Equal(t, expectedError, err, "Should a json error")
}

func TestParseQueryArrayJsonError(t *testing.T) {
	_, err := parseQuery("[[foo]]")

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid character 'o' in literal false (expecting 'a')",
	}
	assert.Equal(t, expectedError, err, "Should a json error")
}

func TestParseQueryHashStatementBasicAnd(t *testing.T) {
	rf, err := parseQuery(`{"logical_operator": "and", "conditions": [["name", "is", "blorg"]]}`)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "and"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	assert.Equal(t, nil, err, "Error should be nil")
	assert.Equal(t, rfExpected, rf, "Should be the same")
}

func TestParseQueryHashStatementBasicOr(t *testing.T) {
	rf, err := parseQuery(`{"logical_operator": "or", "conditions": [["name", "is", "blorg"]]}`)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "or"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	assert.Equal(t, nil, err, "Error should be nil")
	assert.Equal(t, rfExpected, rf, "Should be the same")
}

func TestParseQueryHashStatementAnd(t *testing.T) {
	rf, err := parseQuery(`{"logical_operator": "and", "conditions": [["name", "is", "blorg"], ["sg_status", "in", ["Active", "Bidding"]]]}`)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "and"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))
	rfExpected.AddCondition(newQueryCondition("sg_status", "in", []interface{}{"Active", "Bidding"}))

	assert.Equal(t, nil, err, "Error should be nil")
	assert.Equal(t, rfExpected, rf, "Should be the same")
}

func TestParseQueryHashStatementOr(t *testing.T) {
	rf, err := parseQuery(`{"logical_operator": "or", "conditions": [["name", "is", "blorg"]]}`)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "or"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	assert.Equal(t, nil, err, "Error should be nil")
	assert.Equal(t, rfExpected, rf, "Should be the same")
}

func TestMapToReadFiltersConditionsOnlyError(t *testing.T) {
	rf := newReadFilters()
	query := map[string]interface{}{
		"conditions": []interface{}{[]string{"name", "is", "blorg"}},
	}

	err := mapToReadFilters(query, &rf)

	expectedError := queryParseError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Missing key: 'logical_operator'",
	}
	assert.Equal(t, expectedError, err, "Should a json error")
}

func TestMapToReadFiltersLogicalOpOnlyError(t *testing.T) {
	rf := newReadFilters()
	query := map[string]interface{}{
		"logical_operator": "or",
	}

	err := mapToReadFilters(query, &rf)

	expectedError := queryParseError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Missing key: 'conditions'",
	}
	assert.Equal(t, expectedError, err, "Should a json error")
}
