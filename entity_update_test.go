package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateBadRequestJson(t *testing.T) {
	patchBody := `foo`

	req := patchRequest("/Project/64", patchBody)

	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, `foo`)
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateNoId(t *testing.T) {
	patchBody := `{"name": "My Project"}`
	req := patchRequest("/Shot/", patchBody)

	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, `foo`)
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusNotFound, w.Code)
}

//
func TestUpdateSimple(t *testing.T) {
	patchBody := `{"name": "My Project 2"}`

	req := patchRequest("/Project/75", patchBody)

	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, `{"results":{"id":75,"name":"My Project 2","type":"Project"}}`)
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusOK, w.Code)

	type Entity struct {
		Type string `json:"type"`
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	expectedResponse := Entity{
		Type: "Project",
		Id:   75,
		Name: "My Project 2",
	}

	var jsonResp Entity
	err := json.Unmarshal(w.Body.Bytes(), &jsonResp)
	assert.Nil(t, err)

	assert.Equal(t, expectedResponse, jsonResp)
}

func TestUpdateConflict(t *testing.T) {
	patchBody := `{"code": "SH001"}`

	req := patchRequest("/Shot/2", patchBody)

	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200,
		`{"exception":true,"message":"API update() CRUD ERROR #51: Update failed for [Shot.code]. The value for the Shot Code field is required to be unique. <br> ","error_code":104}`)
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestUpdateBadResponseJson(t *testing.T) {
	patchBody := `{"code": "SH002"}`

	req := patchRequest("/Shot/2", patchBody)

	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, `foo`)
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	assert.Equal(t, http.StatusBadGateway, w.Code)
}
