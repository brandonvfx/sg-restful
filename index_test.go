package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndexSuccess(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200,
		`{"version":[6,0,1],"s3_uploads_enabled":true,"totango_site_id":"0","totango_site_name":"com_shotgunstudio_brandon"}`)

	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
}

func TestIndexShotugnMissing(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(503, "")
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	if w.Code != http.StatusBadGateway {
		t.Errorf("Home page didn't return %v", http.StatusBadGateway)
	}
}

func TestIndexEmptyResponse(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, "")
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Home page didn't return %v", http.StatusInternalServerError)
	}
}

func TestIndexEmptyNonJsonResponse(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server, client, config := mockShotgun(200, "200 Not Ok - Because Shotgun.")
	defer server.Close()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sgConn", *client)
	router(config).ServeHTTP(w, req.WithContext(ctx))
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Home page didn't return %v", http.StatusInternalServerError)
	}
}

func TestIndexNoShotgunClient(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server, _, config := mockShotgun(200, "200 Not Ok - Because Shotgun.")
	defer server.Close()

	router(config).ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Home page didn't return %v", http.StatusInternalServerError)
	}
}

func TestIndexBadShotgunHost(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	ctx := req.Context()
	ctx = context.WithValue(ctx, "sg_conn", Shotgun{
		ServerURL:  "http://localhost:102782/",
		ScriptName: "fake-script",
		ScriptKey:  "fake-key",
		client:     http.Client{},
	})
	config := newClientConfig("0.0.0-test.1", "http://localhost:102782/")

	router(config).ServeHTTP(w, req.WithContext(ctx))
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Home page didn't return %v", http.StatusInternalServerError)
	}
}
