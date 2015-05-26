package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/context"
)

func init() {
	// log.SetLevel(log.PanicLevel)
}

func mockShotgun(code int, body string) (*httptest.Server, *Shotgun) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpClient := &http.Client{Transport: transport}

	server_url, _ := url.Parse(server.URL)
	// Update SG_HOST so that the middleware doesn't break during tesing.
	SG_HOST = server.URL
	var client *Shotgun
	client = &Shotgun{
		ServerUrl:  server_url,
		ScriptName: "fake-script",
		ScriptKey:  "fake-key",
		client:     *httpClient,
	}

	return server, client
}

func TestIndexSuccess(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server, client := mockShotgun(200,
		`{"version":[6,0,1],"s3_uploads_enabled":true,"totango_site_id":"0","totango_site_name":"com_shotgunstudio_brandon"}`)

	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Home page didn't return %v", http.StatusOK)
	}
}

func TestIndexShotugnMissing(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server, client := mockShotgun(503, "")
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)
	if w.Code != http.StatusBadGateway {
		t.Errorf("Home page didn't return %v", http.StatusBadGateway)
	}
}

func TestIndexEmptyResponse(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server, client := mockShotgun(200, "")
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Home page didn't return %v", http.StatusInternalServerError)
	}
}

func TestIndexEmptyNonJsonResponse(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server, client := mockShotgun(200, "200 Not Ok - Because Shotgun.")
	defer server.Close()

	context.Set(req, "sg_conn", *client)
	Router().ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Home page didn't return %v", http.StatusInternalServerError)
	}
}

func TestIndexNoShotgunClient(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	server, _ := mockShotgun(200, "200 Not Ok - Because Shotgun.")
	defer server.Close()

	Router().ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Home page didn't return %v", http.StatusInternalServerError)
	}
}

func TestIndexBadShotgunHost(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	// server, _ := mockShotgun(200, "200 Not Ok - Because Shotgun.")
	// defer server.Close()

	u, _ := url.Parse("http://localhost:102782/")
	context.Set(req, "sg_conn", Shotgun{
		ServerUrl:  u,
		ScriptName: "fake-script",
		ScriptKey:  "fake-key",
		client:     http.Client{},
	})
	Router().ServeHTTP(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Home page didn't return %v", http.StatusInternalServerError)
	}
}
