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
type Format3TestSuite struct {
	suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *Format3TestSuite) SetupSuite() {
	log.Info(" -- Format3 Test Suite --\n")
	manager := GetQPManager()
	manager.ResetActive()
	manager.AddParser("format3", &Format3{})
	manager.SetActiveParsers("format3")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestFormat3TestSuite(t *testing.T) {
	suite.Run(t, new(Format3TestSuite))
}

func (suite *Format3TestSuite) TestParseQueryArrayJsonError() {
	_, err := parseQuery("[[foo]]")

	expectedError := queryParseError{
		StatusCode: http.StatusBadRequest,
		Message:    "invalid character 'o' in literal false (expecting 'a')",
	}
	suite.Equal(expectedError, err, "Should a json error")
}
