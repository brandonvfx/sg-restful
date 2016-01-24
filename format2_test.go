package main

import (
	"net/http"
	"testing"

	log "github.com/Sirupsen/logrus"

	"github.com/stretchr/testify/suite"
)

// Format2TestSuite defines the suite, and absorbs the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type Format2TestSuite struct {
	suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *Format2TestSuite) SetupSuite() {
	log.Info(" -- Format2 Test Suite --\n")
	manager := GetQPManager()
	manager.ResetActive()
	manager.AddParser("format2", &Format2{})
	manager.SetActiveParsers("format2")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestFormat2TestSuite(t *testing.T) {
	suite.Run(t, new(Format2TestSuite))
}

func (suite *Format2TestSuite) TestParseQueryHashJsonError() {
	_, err := parseQuery("{foo}")

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid character 'f' looking for beginning of object key string",
	}
	suite.Equal(expectedError, err, "Should a json error")
}

func (suite *Format2TestSuite) TestParseQueryHashConditionsOnly() {
	_, err := parseQuery(`{"conditions": [["name", "is", "blorg"]]}`)

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "Missing key: 'logical_operator'",
	}
	suite.Equal(expectedError, err, "Should a json error")
}

func (suite *Format2TestSuite) TestParseQueryHashLogicalOpOnly() {
	_, err := parseQuery(`{"logical_operator": "and"}`)

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "Missing key: 'conditions'",
	}
	suite.Equal(expectedError, err, "Should a json error")
}

func (suite *Format2TestSuite) TestParseQueryArrayJsonError() {
	_, err := parseQuery("[[foo]]")

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid character 'o' in literal false (expecting 'a')",
	}
	suite.Equal(expectedError, err, "Should a json error")
}

func (suite *Format2TestSuite) TestParseQueryHashStatementBasicAnd() {
	rf, err := parseQuery(`{"logical_operator": "and", "conditions": [["name", "is", "blorg"]]}`)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "and"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}

func (suite *Format2TestSuite) TestParseQueryHashStatementBasicOr() {
	rf, err := parseQuery(`{"logical_operator": "or", "conditions": [["name", "is", "blorg"]]}`)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "or"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}

func (suite *Format2TestSuite) TestParseQueryHashStatementAnd() {
	rf, err := parseQuery(`{"logical_operator": "and", "conditions": [["name", "is", "blorg"], ["sg_status", "in", ["Active", "Bidding"]]]}`)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "and"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))
	rfExpected.AddCondition(newQueryCondition("sg_status", "in", []interface{}{"Active", "Bidding"}))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}

func (suite *Format2TestSuite) TestParseQueryHashStatementOr() {
	rf, err := parseQuery(`{"logical_operator": "or", "conditions": [["name", "is", "blorg"]]}`)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "or"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}

func (suite *Format2TestSuite) TestMapToReadFiltersConditionsOnlyError() {
	rf := newReadFilters()
	query := map[string]interface{}{
		"conditions": []interface{}{[]string{"name", "is", "blorg"}},
	}

	err := mapToReadFilters(query, &rf)

	expectedError := queryParseError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Missing key: 'logical_operator'",
	}
	suite.Equal(expectedError, err, "Should a json error")
}

func (suite *Format2TestSuite) TestMapToReadFiltersLogicalOpOnlyError() {
	rf := newReadFilters()
	query := map[string]interface{}{
		"logical_operator": "or",
	}

	err := mapToReadFilters(query, &rf)

	expectedError := queryParseError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Missing key: 'conditions'",
	}
	suite.Equal(expectedError, err, "Should a json error")
}
