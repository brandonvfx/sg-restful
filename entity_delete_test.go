package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/context"
	"github.com/stretchr/testify/assert"
)

func deleteRequest(path string) *http.Request {
	req, _ := http.NewRequest("DELETE", path, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", fakeAuthB64))
	return req
}

func TestDeleteNoId(t *testing.T) {
	req := deleteRequest("/Shot/")

	w := httptest.NewRecorder()

	server, client := mockShotgun(200, `foo`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeletePermissionsIssue(t *testing.T) {
	req := deleteRequest("/Project/75")

	w := httptest.NewRecorder()

	server, client := mockShotgun(200,
		`{"exception":true,"message":"API delete() CRUD ERROR #4.1: Entity Project 75 can not be `+
			`deleted by this user. Rule: API Admin -- PermissionRule 315: retire_entity_condition FOR `+
			`entity_type => Project.  RULE: {\"path\":\"name\", \"relation\":\"is_not\",\"values\":`+
			`[\"Template Project\"]}","error_code":104}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDeleteError(t *testing.T) {
	req := deleteRequest("/Project/75")

	w := httptest.NewRecorder()

	server, client := mockShotgun(200,
		`{"exception":true,"message":"Somer Error message","error_code":104}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteSuccess(t *testing.T) {
	req := deleteRequest("/Project/75")

	w := httptest.NewRecorder()

	server, client := mockShotgun(200, `{"results": true}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusGone, w.Code)
}

func TestDeleteMissing(t *testing.T) {
	req := deleteRequest("/Shot/1000000")

	w := httptest.NewRecorder()

	server, client := mockShotgun(200,
		`{"exception":true,"message":"API delete() CRUD ERROR #3: Entity of type [Shot] with id=1000000 does not exist.","error_code":104}`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestDeleteBadJsonResponse(t *testing.T) {
	req := deleteRequest("/Shot/1000000")

	w := httptest.NewRecorder()

	server, client := mockShotgun(200, `foo`)
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadGateway, w.Code)
}
