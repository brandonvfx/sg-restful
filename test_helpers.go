// +build test

package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	log "github.com/Sirupsen/logrus"
)

const fakeAuthB64 = "ZmFrZS1zY3JpcHQ6ZmFrZS1rZXk="

func init() {
	log.SetLevel(log.PanicLevel)
}

func mockShotgun(code int, body string) (*httptest.Server, *Shotgun, clientConfig) {
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

	// Update SG_HOST so that the middleware doesn't break during tesing.
	// SG_HOST = server.URL
	var client *Shotgun
	client = &Shotgun{
		ServerURL:  server.URL,
		ScriptName: "fake-script",
		ScriptKey:  "fake-key",
		client:     *httpClient,
	}

	config := newClientConfig("0.0.0-test.1", server.URL)

	return server, client, config
}

func getRequest(path string) *http.Request {
	req, _ := http.NewRequest("GET", path, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", fakeAuthB64))
	return req
}

func postRequest(path, body string) *http.Request {
	req, _ := http.NewRequest("POST", path, bytes.NewReader([]byte(body)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", fakeAuthB64))
	return req
}

func deleteRequest(path string) *http.Request {
	req, _ := http.NewRequest("DELETE", path, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", fakeAuthB64))
	return req
}

func patchRequest(path, body string) *http.Request {
	req, _ := http.NewRequest("PATCH", path, bytes.NewReader([]byte(body)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", fakeAuthB64))
	return req
}
