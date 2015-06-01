package main

import (
	"encoding/json"
	"fmt"

	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/context"
)

const fakeAuthB64 = "ZmFrZS1zY3JpcHQ6ZmFrZS1rZXk="

func init() {
	//log.SetLevel(log.PanicLevel)
}

func getRequest(path string) *http.Request {
	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", fakeAuthB64))
	return req
}

func TestReadFindOneSimple(t *testing.T) {
	req := getRequest("/Project/65")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client := mockShotgun(200,
		`{"results":{"entities":[{"type":"Project","id":65}],"paging_info":{"current_page":1,"page_count":1,"entity_count":1,"entities_per_page":1}}}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Project read did not return %v", http.StatusOK)
	}

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
	}

	expected := `{"type":"Project","id":65}`
	var jsonExpected Entity
	err := json.Unmarshal(w.Body.Bytes(), &jsonExpected)
	if err != nil {
		t.Errorf("Issue Unmarshaling Excpected")
	}

	var jsonResp Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Errorf("Issue Unmarshaling Response")
	}

	if jsonResp != jsonExpected {
		t.Errorf("response didn't match. resp: %v != expected %v", jsonResp, expected)
	}
}

func TestReadFindOneWithFields(t *testing.T) {
	req := getRequest("/Project/65?fields=name,sg_status")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client := mockShotgun(200,
		`{"results":{"entities":[{"id":65,"name":"Big Buck Bunny","sg_status":"Active","type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":1,"entities_per_page":1}}}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Project read did not return %v", http.StatusOK)
	}

	type Entity struct {
		Type   string `json:"type"`
		Id     int    `json:"id"`
		Name   string `json:"name"`
		Status string `json:"sg_status"`
	}

	expected := `{"id":65,"name":"Big Buck Bunny","sg_status":"Active","type":"Project"}`
	var jsonExpected Entity
	err := json.Unmarshal(w.Body.Bytes(), &jsonExpected)
	if err != nil {
		t.Errorf("Issue Unmarshaling Excpected")
	}

	var jsonResp Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Errorf("Issue Unmarshaling Response")
	}

	if jsonResp != jsonExpected {
		t.Errorf("response didn't match. resp: %v != expected %v", jsonResp, expected)
	}
}

func TestReadFindOneNoneIntId(t *testing.T) {
	req := getRequest("/Project/foo")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client := mockShotgun(200,
		`{"results":{"entities":[{"id":65, "type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":1,"entities_per_page":1}}}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Project read did not return %v", http.StatusNotFound)
	}
}

func TestReadNoShotgunClient(t *testing.T) {
	req := getRequest("/Project")
	w := httptest.NewRecorder()

	server, _ := mockShotgun(200, "200 Not Ok - Because Shotgun.")
	defer server.Close()

	Router().ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Home page didn't return %v", http.StatusInternalServerError)
	}
}
