package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
)

func TestFindAllSimple(t *testing.T) {
	req := getRequest("/Project")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client := mockShotgun(200,
		`{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
	}

	expected := `[{"type":"Project","id":63},{"type":"Project","id":65}]`
	var jsonExpected []Entity
	err := json.Unmarshal([]byte(expected), &jsonExpected)
	if err != nil {
		t.Errorf("Issue Unmarshaling Excpected")
	}
	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Errorf("Issue Unmarshaling Response")
	}

	assert.Equal(t, jsonExpected, jsonResp)
}

func TestFindAllWithFields(t *testing.T) {
	req := getRequest("/Project?fields=name,sg_status")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client := mockShotgun(200,
		`{"results":{"entities":[{"id":63,"name":"Template Project","sg_status":null,"type":"Project"},{"id":65,"name":"Big Buck Bunny","sg_status":"Active","type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	type Entity struct {
		Type   string `json:"type"`
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Status string `json:"sg_status"`
	}

	expected := `[{"id":63,"name":"Template Project","sg_status":null,"type":"Project"},{"id":65,"name":"Big Buck Bunny","sg_status":"Active","type":"Project"}]`
	var jsonExpected []Entity
	err := json.Unmarshal([]byte(expected), &jsonExpected)
	if err != nil {
		t.Errorf("Issue Unmarshaling Excpected")
	}
	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Errorf("Issue Unmarshaling Response")
	}

	assert.Equal(t, jsonExpected, jsonResp)
}

func TestFindAllPagingLimit1(t *testing.T) {
	req := getRequest("/Project?limit=1&page=2")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client := mockShotgun(200,
		`{"results":{"entities":[{"id":63,"name":"Template Project","sg_status":null,"type":"Project"},{"id":65,"name":"Big Buck Bunny","sg_status":"Active","type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
	}

	expected := `[{"type":"Project","id":63}]`
	var jsonExpected []Entity
	err := json.Unmarshal([]byte(expected), &jsonExpected)
	if err != nil {
		t.Errorf("Issue Unmarshaling Excpected")
	}
	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Errorf("Issue Unmarshaling Response")
	}

}

func TestFindAllBadPageValue(t *testing.T) {
	req := getRequest("/Project?limit=1&page=foo")
	w := httptest.NewRecorder()

	server, client := mockShotgun(200,
		`{"exception":true,"message":"API read() invalid/missing integer 'paging' 'entities_per_page':\n{\"current_page\"=>\"foo\", \"entities_per_page\"=>1}","error_code":103}`)
	defer server.Close()
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestFindAllBadLimitValue(t *testing.T) {
	req := getRequest("/Project?limit=foo&page=1")
	w := httptest.NewRecorder()

	server, client := mockShotgun(200,
		`{"exception":true,"message":"API read() invalid/missing integer 'paging' 'entities_per_page':\n{\"current_page\"=>2, \"entities_per_page\"=>\"foo\"}","error_code":103}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Coming soon
func TestFindAllNoResults(t *testing.T) {
	req := getRequest(`/Project?q=[["name", "is", "foo"]]`)
	w := httptest.NewRecorder()

	server, client := mockShotgun(200,
		`{"results":{"entities":[],"paging_info":{"current_page":0,"page_count":0,"entity_count":0,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestFindAllMapQuery(t *testing.T) {
	req := getRequest(`/Project?q={"logical_operator": "and", "conditions": [["name", "starts_with", "Big"]]}`)
	w := httptest.NewRecorder()

	server, client := mockShotgun(200,
		`{"results":{"entities":[{"id":65, "type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":1,"entities_per_page":1}}}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
	}

	expected := `[{"type":"Project","id":65}]`
	var jsonExpected []Entity
	err := json.Unmarshal([]byte(expected), &jsonExpected)
	if err != nil {
		t.Errorf("Issue Unmarshaling Excpected")
	}
	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Errorf("Issue Unmarshaling Response")
	}

	assert.Equal(t, jsonExpected, jsonResp)

}

func TestFindAllArrayQuery(t *testing.T) {
	req := getRequest(`/Project?q=[["name", "starts_with", "Big"]]`)
	w := httptest.NewRecorder()

	server, client := mockShotgun(200,
		`{"results":{"entities":[{"id":65, "type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":1,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
	}

	expected := `[{"type":"Project","id":65}]`
	var jsonExpected []Entity
	err := json.Unmarshal([]byte(expected), &jsonExpected)
	if err != nil {
		t.Errorf("Issue Unmarshaling Excpected")
	}
	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Errorf("Issue Unmarshaling Response")
	}

	assert.Equal(t, jsonExpected, jsonResp)

}

func TestFindAllAndStringQuery(t *testing.T) {
	req := getRequest(`/Project?q=and(["name", "starts_with", "Big"])`)
	w := httptest.NewRecorder()

	server, client := mockShotgun(200,
		`{"results":{"entities":[{"id":65, "type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":1,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
	}

	expected := `[{"type":"Project","id":65}]`
	var jsonExpected []Entity
	err := json.Unmarshal([]byte(expected), &jsonExpected)
	if err != nil {
		t.Errorf("Issue Unmarshaling Excpected")
	}
	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Errorf("Issue Unmarshaling Response")
	}

	assert.Equal(t, jsonExpected, jsonResp)

}

func TestFindAllOrStringQuery(t *testing.T) {
	req := getRequest(`/Project?q=or(["name", "starts_with", "Big"],["name", "starts_with", "Test"])`)
	w := httptest.NewRecorder()

	server, client := mockShotgun(200,
		`{"results":{"entities":[{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":3,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
	}

	expected := `[{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}]`
	var jsonExpected []Entity
	err := json.Unmarshal([]byte(expected), &jsonExpected)
	if err != nil {
		t.Errorf("Issue Unmarshaling Excpected")
	}
	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Errorf("Issue Unmarshaling Response")
	}

	assert.Equal(t, jsonExpected, jsonResp)

}

func TestFindAllBadQuery(t *testing.T) {
	req := getRequest(`/Project?q=and(["name", "starts_with", "Test"],)`)
	w := httptest.NewRecorder()

	server, client := mockShotgun(200,
		`{"results":{"entities":[{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":3,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

}

func TestFindAllBadAuth(t *testing.T) {
	req, _ := http.NewRequest("GET", `/Project`, nil)
	w := httptest.NewRecorder()

	server, _ := mockShotgun(200,
		`{"exception":true,"message":"Can't authenticate script 'TestScript'","error_code":102}`)
	defer server.Close()

	Router().ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)

}
