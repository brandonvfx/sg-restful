package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/Sirupsen/logrus"

	"github.com/gorilla/context"
	"github.com/stretchr/testify/suite"
)

// EntityGetAllTestSuite defines the suite, and absorbs the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type EntityGetAllTestSuite struct {
	suite.Suite
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *EntityGetAllTestSuite) SetupSuite() {
	manager := GetQPManager()
	log.Info(" -- EntityGetAll Test Suite --\n")
	log.Debug("EntityGetAllTestSuite.SetupSuite() - setting active parsers to format1, format2, format3")
	manager.SetActiveParsers("format1", "format2", "format3")
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestEntityGetAllTestSuite(t *testing.T) {
	log.Debug("EntityGetALlTestSuite - Running test suite")
	suite.Run(t, new(EntityGetAllTestSuite))
	log.Debug("EntityGetAllTestSuite - Finished test suite")
}

func (suite *EntityGetAllTestSuite) TestFindAllSimple() {
	req := getRequest("/Project")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client, config := mockShotgun(200,
		`{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
	}

	expected := `[{"type":"Project","id":63},{"type":"Project","id":65}]`
	var jsonExpected []Entity
	err := json.Unmarshal([]byte(expected), &jsonExpected)
	suite.Nil(err)

	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	suite.Nil(err)

	suite.Equal(jsonExpected, jsonResp)
}

func (suite *EntityGetAllTestSuite) TestFindAllWithFields() {
	req := getRequest("/Project?fields=name,sg_status")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client, config := mockShotgun(200,
		`{"results":{"entities":[{"id":63,"name":"Template Project","sg_status":null,"type":"Project"},{"id":65,"name":"Big Buck Bunny","sg_status":"Active","type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	type Entity struct {
		Type   string `json:"type"`
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Status string `json:"sg_status"`
	}

	expected := `[{"id":63,"name":"Template Project","sg_status":null,"type":"Project"},{"id":65,"name":"Big Buck Bunny","sg_status":"Active","type":"Project"}]`
	var jsonExpected []Entity
	err := json.Unmarshal([]byte(expected), &jsonExpected)
	suite.Nil(err)

	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	suite.Nil(err)

	suite.Equal(jsonExpected, jsonResp)
}

func (suite *EntityGetAllTestSuite) TestFindAllPagingLimit1() {
	req := getRequest("/Project?limit=1&page=2")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client, config := mockShotgun(200,
		`{"results":{"entities":[{"id":63,"name":"Template Project","sg_status":null,"type":"Project"},{"id":65,"name":"Big Buck Bunny","sg_status":"Active","type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
	}

	expected := `[{"type":"Project","id":63}]`
	var jsonExpected []Entity
	err := json.Unmarshal([]byte(expected), &jsonExpected)
	suite.Nil(err)

	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	suite.Nil(err)

}

func (suite *EntityGetAllTestSuite) TestFindAllBadPageValue() {
	req := getRequest("/Project?limit=1&page=foo")
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200,
		`{"exception":true,"message":"API read() invalid/missing integer 'paging' 'entities_per_page':\n{\"current_page\"=>\"foo\", \"entities_per_page\"=>1}","error_code":103}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)

}

func (suite *EntityGetAllTestSuite) TestFindAllBadLimitValue() {
	req := getRequest("/Project?limit=foo&page=1")
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200,
		`{"exception":true,"message":"API read() invalid/missing integer 'paging' 'entities_per_page':\n{\"current_page\"=>2, \"entities_per_page\"=>\"foo\"}","error_code":103}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)
}

// Coming soon
func (suite *EntityGetAllTestSuite) TestFindAllNoResults() {
	req := getRequest(`/Project?q=[["name", "is", "foo"]]`)
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200,
		`{"results":{"entities":[],"paging_info":{"current_page":0,"page_count":0,"entity_count":0,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)

	suite.Equal(http.StatusNoContent, w.Code)
}

func (suite *EntityGetAllTestSuite) TestFindAllMapQuery() {
	req := getRequest(`/Project?q={"logical_operator": "and", "conditions": [["name", "starts_with", "Big"]]}`)
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200,
		`{"results":{"entities":[{"id":65, "type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":1,"entities_per_page":1}}}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
	}

	expected := `[{"type":"Project","id":65}]`
	var jsonExpected []Entity
	err := json.Unmarshal([]byte(expected), &jsonExpected)
	suite.Nil(err)

	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	suite.Nil(err)

	suite.Equal(jsonExpected, jsonResp)

}

func (suite *EntityGetAllTestSuite) TestFindAllArrayQuery() {
	req := getRequest(`/Project?q=[["name", "starts_with", "Big"]]`)
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200,
		`{"results":{"entities":[{"id":65, "type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":1,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
	}

	expected := `[{"type":"Project","id":65}]`
	var jsonExpected []Entity
	err := json.Unmarshal([]byte(expected), &jsonExpected)
	suite.Nil(err)

	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	suite.Nil(err)

	suite.Equal(jsonExpected, jsonResp)

}

func (suite *EntityGetAllTestSuite) TestFindAllAndStringQuery() {
	req := getRequest(`/Project?q=and(["name", "starts_with", "Big"])`)
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200,
		`{"results":{"entities":[{"id":65, "type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":1,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
	}

	expected := `[{"type":"Project","id":65}]`
	var jsonExpected []Entity
	err := json.Unmarshal([]byte(expected), &jsonExpected)
	suite.Nil(err)

	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	suite.Nil(err)

	suite.Equal(jsonExpected, jsonResp)

}

func (suite *EntityGetAllTestSuite) TestFindAllOrStringQuery() {
	req := getRequest(`/Project?q=or(["name", "starts_with", "Big"],["name", "starts_with", "Test"])`)
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200,
		`{"results":{"entities":[{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":3,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)

	suite.Equal(http.StatusOK, w.Code)

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
	}

	expected := `[{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}]`
	var jsonExpected []Entity
	err := json.Unmarshal([]byte(expected), &jsonExpected)
	suite.Nil(err)

	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	suite.Nil(err)

	suite.Equal(jsonExpected, jsonResp)

}

func (suite *EntityGetAllTestSuite) TestFindAllBadQuery() {
	req := getRequest(`/Project?q=and(["name", "starts_with", "Test"],)`)
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200,
		`{"results":{"entities":[{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":3,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)

}

func (suite *EntityGetAllTestSuite) TestFindAllBadAuth() {
	req, _ := http.NewRequest("GET", `/Project`, nil)
	w := httptest.NewRecorder()

	server, _, config := mockShotgun(200,
		`{"exception":true,"message":"Can't authenticate script 'TestScript'","error_code":102}`)
	defer server.Close()

	router(config).ServeHTTP(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
}

func (suite *EntityGetAllTestSuite) TestFindAllBadResponseJson() {
	req := getRequest("/Project")
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, `foo`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)
	suite.Equal(http.StatusBadGateway, w.Code)
}

func (suite *EntityGetAllTestSuite) TestFindAllAnException() {
	req := getRequest("/Project")
	w := httptest.NewRecorder()

	server, _, config := mockShotgun(200,
		`{"exception":true,"message":"Can't authenticate script 'TestScript'","error_code":102}`)
	defer server.Close()

	router(config).ServeHTTP(w, req)

	suite.Equal(http.StatusInternalServerError, w.Code)
}
