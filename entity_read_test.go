package main

import (
	"encoding/json"
	"fmt"

	"net/http"
	"net/http/httptest"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
)

const fakeAuthB64 = "ZmFrZS1zY3JpcHQ6ZmFrZS1rZXk="

func init() {
	log.SetLevel(log.PanicLevel)
}

func getRequest(method, path string) *http.Request {
	req, _ := http.NewRequest(method, path, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", fakeAuthB64))
	return req
}

func TestReadFindAllSimple(t *testing.T) {
	req := getRequest("GET", "/Project")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client := mockShotgun(200,
		`{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)
	log.Println(w.Code)
	log.Println(w.Body.String())
	if w.Code != http.StatusOK {
		t.Errorf("Project read did not return %v", http.StatusOK)
	}

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
	}

	expected := `[{"type":"Project","id":63},{"type":"Project","id":65}]`
	var jsonExpected []Entity
	err := json.Unmarshal(w.Body.Bytes(), &jsonExpected)
	if err != nil {
		t.Errorf("Issue Unmarshaling Excpected")
	}
	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Errorf("Issue Unmarshaling Response")
	}

	if len(jsonResp) != len(jsonExpected) {
		t.Errorf("Respnse size didn't match. resp: %v expected %v", jsonResp, expected)
	}

	for i := range jsonResp {
		if jsonResp[i] != jsonExpected[i] {
			t.Errorf("response didn't match, item %d. resp: %v != expected %v", i, jsonResp[i], expected[i])
		}
	}
}

func TestReadFindAllWithFields(t *testing.T) {
	req := getRequest("GET", "/Project?fields=name,sg_status")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client := mockShotgun(200,
		`{"results":{"entities":[{"id":63,"name":"Template Project","sg_status":null,"type":"Project"},{"id":65,"name":"Big Buck Bunny","sg_status":"Active","type":"Project"}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}`)
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

	expected := `[{"id":63,"name":"Template Project","sg_status":null,"type":"Project"},{"id":65,"name":"Big Buck Bunny","sg_status":"Active","type":"Project"}]`
	var jsonExpected []Entity
	err := json.Unmarshal(w.Body.Bytes(), &jsonExpected)
	if err != nil {
		t.Errorf("Issue Unmarshaling Excpected")
	}
	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Errorf("Issue Unmarshaling Response")
	}

	if len(jsonResp) != len(jsonExpected) {
		t.Errorf("Respnse size didn't match. resp: %v expected %v", jsonResp, expected)
	}

	for i := range jsonResp {
		if jsonResp[i] != jsonExpected[i] {
			t.Errorf("response didn't match, item %d. resp: %v != expected %v", i, jsonResp[i], expected[i])
		}
	}
}

func TestReadFindAllPagingLimit1(t *testing.T) {
	req := getRequest("GET", "/Project?limit=1&page=2")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client := mockShotgun(200,
		`{"results":{"entities":[{"type":"Project","id":65}],"paging_info":{"current_page":2,"page_count":4,"entity_count":4,"entities_per_page":1}}}`)
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

	expected := `[{"type":"Project","id":65}]`
	var jsonExpected []Entity
	err := json.Unmarshal(w.Body.Bytes(), &jsonExpected)
	if err != nil {
		t.Errorf("Issue Unmarshaling Excpected")
	}
	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Errorf("Issue Unmarshaling Response")
	}

	if len(jsonResp) != len(jsonExpected) {
		t.Errorf("Respnse size didn't match. resp: %v expected %v", jsonResp, expected)
	}

	for i := range jsonResp {
		if jsonResp[i] != jsonExpected[i] {
			t.Errorf("response didn't match, item %d. resp: %v != expected %v", i, jsonResp[i], expected[i])
		}
	}
}

func TestReadFindAllBadPageValue(t *testing.T) {
	req := getRequest("GET", "/Project?limit=1&page=foo")
	w := httptest.NewRecorder()

	//{"results":{"entities":[{"type":"Project","id":63},{"type":"Project","id":65},{"type":"Project","id":66},{"type":"Project","id":71}],"paging_info":{"current_page":1,"page_count":1,"entity_count":4,"entities_per_page":500}}}
	server, client := mockShotgun(200,
		`{"results":{"entities":[{"type":"Project","id":63}],"paging_info":{"current_page":2,"page_count":4,"entity_count":4,"entities_per_page":1}}}`)
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

	expected := `[{"type":"Project","id":63}]`
	var jsonExpected []Entity
	err := json.Unmarshal(w.Body.Bytes(), &jsonExpected)
	if err != nil {
		t.Errorf("Issue Unmarshaling Excpected")
	}
	var jsonResp []Entity
	err = json.Unmarshal(w.Body.Bytes(), &jsonResp)
	if err != nil {
		t.Errorf("Issue Unmarshaling Response")
	}

	if len(jsonResp) != len(jsonExpected) {
		t.Errorf("Respnse size didn't match. resp: %v expected %v", jsonResp, expected)
	}

	for i := range jsonResp {
		if jsonResp[i] != jsonExpected[i] {
			t.Errorf("response didn't match, item %d. resp: %v != expected %v", i, jsonResp[i], expected[i])
		}
	}
}

// Coming soon
// func TestReadFindAllNoResults(t *testing.T) {
// 	req, _ := http.NewRequest("GET", "/Project?name=foo", nil)
// 	w := httptest.NewRecorder()

// 	server, client := mockShotgun(204, "")
// 	defer server.Close()

// 	context.Set(req, "sg_conn", *client)
// 	Router().ServeHTTP(w, req)
// 	if w.Code != http.StatusNoContent {
// 		t.Errorf("Project read did not return %v", http.StatusNoContent)
// 	}
// }

func TestReadFindOneSimple(t *testing.T) {
	req := getRequest("GET", "/Project/65")
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
	req := getRequest("GET", "/Project/65?fields=name,sg_status")
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
	req := getRequest("GET", "/Project/foo")
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
	req := getRequest("GET", "/Project")
	w := httptest.NewRecorder()

	server, _ := mockShotgun(200, "200 Not Ok - Because Shotgun.")
	defer server.Close()

	Router().ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Home page didn't return %v", http.StatusInternalServerError)
	}
}
