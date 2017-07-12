package main

import (
	"net/http"
	"testing"

	log "github.com/sirupsen/logrus"

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
	f := &Format2{}
	testString := "{foo}"
	tf := f.CanParseString(testString)
	suite.Equal(false, tf, "Should not be able to parse string")

	_, err := f.ParseString(testString)

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid character 'f' looking for beginning of object key string",
	}
	suite.Equal(expectedError, err, "Should a json error")
}

func (suite *Format2TestSuite) TestParseQueryHashConditionsOnly() {
	f := &Format2{}
	teststring := `{"conditions": [["name", "is", "blorg"]]}`
	tf := f.CanParseString(teststring)
	suite.Equal(tf, false, "Should not be able to parse string")

	_, err := f.ParseString(teststring)

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "Missing key: 'logical_operator'",
	}
	suite.Equal(expectedError, err, "Should a json error")
}

func (suite *Format2TestSuite) TestParseQueryHashLogicalOpOnly() {
	f := &Format2{}
	testString := `{"logical_operator": "and"}`
	tf := f.CanParseString(testString)
	suite.Equal(false, tf, "Should not be able to parse string")

	_, err := f.ParseString(testString)

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "Missing key: 'conditions'",
	}
	suite.Equal(expectedError, err, "Should a json error")
}

func (suite *Format2TestSuite) TestParseQueryHashStatementBasicAnd() {
	f := &Format2{}
	testString := `{"logical_operator": "and", "conditions": [["name", "is", "blorg"]]}`
	tf := f.CanParseString(testString)
	suite.Equal(true, tf, "Should be able to parse string")

	rf, err := f.ParseString(testString)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "and"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}

func (suite *Format2TestSuite) TestParseQueryHashStatementBasicOr() {
	f := &Format2{}
	testString := `{"logical_operator": "or", "conditions": [["name", "is", "blorg"]]}`
	tf := f.CanParseString(testString)
	suite.Equal(true, tf, "Should be able to parse string")

	rf, err := f.ParseString(testString)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "or"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}

func (suite *Format2TestSuite) TestParseQueryHashStatementAnd() {
	testString := `{"logical_operator": "and", "conditions": [["name", "is", "blorg"], ["sg_status", "in", ["Active", "Bidding"]]]}`
	f := &Format2{}

	tf := f.CanParseString(testString)
	suite.Equal(true, tf, "Should be able to parse string")

	rf, err := f.ParseString(testString)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "and"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))
	rfExpected.AddCondition(newQueryCondition("sg_status", "in", []interface{}{"Active", "Bidding"}))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}

func (suite *Format2TestSuite) TestParseQueryHashStatementOr() {
	testString := `{"logical_operator": "or", "conditions": [["name", "is", "blorg"]]}`
	f := &Format2{}

	tf := f.CanParseString(testString)
	suite.Equal(true, tf, "Should be able to parse string")

	rf, err := f.ParseString(testString)

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
