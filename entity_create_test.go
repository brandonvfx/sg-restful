package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateBadRequestJson(t *testing.T) {
	postBody := `foo`

	req := postRequest("/Project", postBody)

	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, `foo`)
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateSimple(t *testing.T) {
	postBody := `{"name": "My Project"}`

	req := postRequest("/Project", postBody)

	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, `{"results":{"id":75,"name":"My Project","type":"Project"}}`)
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusCreated, w.Code)

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	expectedResponse := Entity{
		Type: "Project",
		Id:   75,
		Name: "My Project",
	}

	var jsonResp Entity
	err := json.Unmarshal(w.Body.Bytes(), &jsonResp)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse, jsonResp)
}

//{"exception":true,"message":"API create() CRUD ERROR #61: Create failed for [Project]. The value for the Project Name field is required to be unique. <br>","error_code":104}

func TestCreateConflict(t *testing.T) {
	postBody := `{"name": "My Project"}`

	req := postRequest("/Project", postBody)

	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200,
		`{"exception":true,"message":"API create() CRUD ERROR #61: Create failed for [Project]. The value for the Project Name field is required to be unique. <br>","error_code":104}`)
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestCreateBadResponseJson(t *testing.T) {
	postBody := `{"name": "My Project"}`

	req := postRequest("/Project", postBody)

	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, `foo`)
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusBadGateway, w.Code)
}

func TestCreateShotgunError(t *testing.T) {
	postBody := `{"name": "My Project"}`

	req := postRequest("/Project", postBody)

	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200,
		`{"exception":true,"message":"API create() CRUD ERROR Some other error","error_code":104}`)
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
