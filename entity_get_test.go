package main

import (
	"context"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type EntityGetTestSuite struct {
	suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *EntityGetTestSuite) SetupSuite() {
	manager := GetQPManager()
	log.Info(" -- EntityGet Test Suite --\n")
	log.Debug("EntityGetTestSuite.SetupSuite() - setting active parsers to format1, format2, format3")
	manager.SetActiveParsers("format1", "format2", "format3")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestEntityGetTestSuite(t *testing.T) {
	//ts := new(Format1TestSuite)
	log.Debug("EntityGetTestSuite - Running test suite")
	suite.Run(t, new(EntityGetTestSuite))
	log.Debug("EntityGetTestSuite - Finished test suite")
}

func (suite *EntityGetTestSuite) TestReadFindOneSimple() {
	req := getRequest("/Project/65")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client, config := mockShotgun(200,
		`{"results":{"entities":[{"type":"Project","id":65}],"paging_info":{"current_page":1,"page_count":1,"entity_count":1,"entities_per_page":1}}}`)
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	suite.Equal(http.StatusOK, w.Code)

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
	}

	expected := []byte(`{"type":"Project","id":65}`)
	var jsonExpected Entity
	err := json.Unmarshal(expected, &jsonExpected)
	suite.Nil(err)

	var jsonResp Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	suite.Nil(err)

	suite.Equal(jsonExpected, jsonResp)
}

func (suite *EntityGetTestSuite) TestReadFindOneWithFields() {
	req := getRequest("/Project/65?fields=name,sg_status")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client, config := mockShotgun(200,
		`{"results":{"entities":[{"id":65,"name":"Big Buck Bunny","sg_status":"Active","type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":1,"entities_per_page":1}}}`)
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	suite.Equal(http.StatusOK, w.Code)

	type Entity struct {
		Type   string `json:"type"`
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Status string `json:"sg_status"`
	}

	expected := []byte(`{"id":65,"name":"Big Buck Bunny","sg_status":"Active","type":"Project"}`)
	var jsonExpected Entity
	err := json.Unmarshal(expected, &jsonExpected)
	suite.Nil(err)

	var jsonResp Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	suite.Nil(err)

	suite.Equal(jsonExpected, jsonResp)
}

func (suite *EntityGetTestSuite) TestReadFindOneNoneIntId() {
	req := getRequest("/Project/foo")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client, config := mockShotgun(200,
		`{"results":{"entities":[{"id":65, "type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":1,"entities_per_page":1}}}`)
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))

	suite.Equal(http.StatusNotFound, w.Code)
}

func (suite *EntityGetTestSuite) TestGetNotFound() {
	req := getRequest("/Project/64")
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, `{"results":{"entities":[]}}`)
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	suite.Equal(http.StatusNotFound, w.Code)
}

func (suite *EntityGetTestSuite) TestReadNoShotgunClient() {
	req := getRequest("/Project")
	w := httptest.NewRecorder()

	server, _, config := mockShotgun(200, "200 Not Ok - Because Shotgun.")
	defer server.Close()

	router(config).ServeHTTP(w, req)
	suite.Equal(http.StatusBadGateway, w.Code)
}

func (suite *EntityGetTestSuite) TestGetBadResponseJson() {
	req := getRequest("/Project/64")
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, `foo`)
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	suite.Equal(http.StatusBadGateway, w.Code)
}
