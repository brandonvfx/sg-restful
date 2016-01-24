package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"

	"github.com/stretchr/testify/suite"
)

func TestMain(m *testing.M) {
	//f := SetupLog()
	log.Info("TestMain setup")
	manager := GetQPManager()
	res := m.Run()
	f.Close()
	os.Exit(res)
}

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type Format1TestSuite struct {
	suite.Suite
}

func SetupLog() *os.File {
	f, err := os.OpenFile("sg-restful-test.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}

	log.SetLevel(log.DebugLevel)
	log.SetOutput(f)
	log.SetFormatter(&log.TextFormatter{})
	return f

}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *Format1TestSuite) SetupSuite() {
	log.Info(" -- Format1 Test Suite --\n")
	log.Debug("Format1TestSuite.SetupSuite() - setting active parsers to format1")
	manager.SetActiveParsers("format1")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestFormat1TestSuite(t *testing.T) {
	//ts := new(Format1TestSuite)
	log.Info("TestFormat1TestSuite - Running test suite")
	suite.Run(t, new(Format1TestSuite))
	log.Info("TestFormat1TestSuite - Finished test suite")

}

func (suite *Format1TestSuite) TestActiveParsers() {
	log.Info("Format1TestSuite.TestActiveParsers()")
	manager := GetQPManager()
	keys, active := manager.GetActiveParsers()
	log.Infof("active parsers:%s %s", keys, active)
	suite.Equal(len(active), 1, "Should be a single active parser")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementRegExpError() {
	_, err := parseQuery("()")

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "No QueryFormats can parse input",
	}
	suite.Equal(expectedError, err, "Should be formating error")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementFilterFormatError() {
	_, err := parseQuery("and(foo, bar)")

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "No QueryFormats can parse input",
	}
	suite.Equal(expectedError, err, "Should be formating error")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementJsonError() {
	_, err := parseQuery("and([foo, bar,])")

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid character 'o' in literal false (expecting 'a')",
	}
	suite.Equal(expectedError, err, "Should be formating error")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementNotAQuery() {
	_, err := parseQuery("andwell_this_is_bad")

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "No QueryFormats can parse input",
	}
	suite.Equal(expectedError, err, "Should be formating error")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementSliceOfFilters() {
	rf, err := parseQuery(`and([["name", "is", "blorg"]])`)

	rfExpected := newReadFilters()
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementBasicAnd() {

	rf, err := parseQuery(`and(["name", "is", "blorg"])`)

	rfExpected := newReadFilters()
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementBasicAndUpper() {
	rf, err := parseQuery(`AND(["name", "is", "blorg"])`)

	rfExpected := newReadFilters()
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementBasicOr() {
	rf, err := parseQuery(`or(["name", "is", "blorg"])`)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "or"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}

func (suite *Format1TestSuite) TestParseQueryAndOrStatementBasicOrUpper() {
	rf, err := parseQuery(`OR(["name", "is", "blorg"])`)

	rfExpected := newReadFilters()
	rfExpected.LogicalOperator = "or"
	rfExpected.AddCondition(newQueryCondition("name", "is", "blorg"))

	suite.Equal(nil, err, "Error should be nil")
	suite.Equal(rfExpected, rf, "Should be the same")
}
