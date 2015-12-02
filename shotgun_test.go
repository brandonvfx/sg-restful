package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShotgunScript(t *testing.T) {
	host := "http://localhost"
	script := "fake-script"
	key := "fake-key"

	sg := NewShotgun(host, script, key)

	expectedCreds := map[string]string{
		"script_name": "fake-script",
		"script_key":  "fake-key",
	}
	creds := sg.Creds()

	assert.Equal(t, expectedCreds, creds)
	assert.NotEmpty(t, sg.ServerURL)
}

func TestShotgunUser(t *testing.T) {
	host := "http://localhost"
	login := "fake-login"
	password := "fake-pass"

	sg := NewUserShotgun(host, login, password)

	expectedCreds := map[string]string{
		"user_login":    "fake-login",
		"user_password": "fake-pass",
	}

	creds := sg.Creds()
	assert.Equal(t, expectedCreds, creds)
	assert.NotEmpty(t, sg.ServerURL)

}

func TestShotgunGetFullUrl(t *testing.T) {
	fullURL := getFullURL("http://localhost")
	assert.Equal(t, "http://localhost/api3/json", fullURL)
}
