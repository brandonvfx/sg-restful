package main

import (
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
)

func TestReadFindOneSimple(t *testing.T) {
	req := getRequest("/Project/65")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client, config := mockShotgun(200,
		`{"results":{"entities":[{"type":"Project","id":65}],"paging_info":{"current_page":1,"page_count":1,"entity_count":1,"entities_per_page":1}}}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
	}

	expected := []byte(`{"type":"Project","id":65}`)
	var jsonExpected Entity
	err := json.Unmarshal(expected, &jsonExpected)
	assert.Nil(t, err)

	var jsonResp Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	assert.Nil(t, err)

	assert.Equal(t, jsonExpected, jsonResp)
}

func TestReadFindOneWithFields(t *testing.T) {
	req := getRequest("/Project/65?fields=name,sg_status")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client, config := mockShotgun(200,
		`{"results":{"entities":[{"id":65,"name":"Big Buck Bunny","sg_status":"Active","type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":1,"entities_per_page":1}}}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)

	type Entity struct {
		Type   string `json:"type"`
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Status string `json:"sg_status"`
	}

	expected := []byte(`{"id":65,"name":"Big Buck Bunny","sg_status":"Active","type":"Project"}`)
	var jsonExpected Entity
	err := json.Unmarshal(expected, &jsonExpected)
	assert.Nil(t, err)

	var jsonResp Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	assert.Nil(t, err)

	assert.Equal(t, jsonExpected, jsonResp)
}

func TestReadFindOneNoneIntId(t *testing.T) {
	req := getRequest("/Project/foo")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client, config := mockShotgun(200,
		`{"results":{"entities":[{"id":65, "type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":1,"entities_per_page":1}}}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetNotFound(t *testing.T) {
	req := getRequest("/Project/64")
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, `{"results":{"entities":[]}}`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestReadNoShotgunClient(t *testing.T) {
	req := getRequest("/Project")
	w := httptest.NewRecorder()

	server, _, config := mockShotgun(200, "200 Not Ok - Because Shotgun.")
	defer server.Close()

	router(config).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadGateway, w.Code)
}

func TestGetBadResponseJson(t *testing.T) {
	req := getRequest("/Project/64")
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, `foo`)
	defer server.Close()

	context.Set(req, "sgConn", *client)
	router(config).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadGateway, w.Code)
}
