package main

import "testing"

func TestShotgunScript(t *testing.T) {
	host := "http://localhost"
	script := "fake-script"
	key := "fake-key"
	sg := NewShotgun(
		host,
		script,
		key,
	)

	expectedUrl := host + "/api3/json"
	if sg.ServerUrl.String() != expectedUrl {
		t.Errorf("%v != %s", sg.ServerUrl.String())
	}

	creds := sg.Creds()

	credsScript, ok := creds["script_name"]
	if !ok {
		t.Error("Creds map missing script_name")
	}
	if credsScript != script {
		t.Errorf("Script name miss-match: %v != %v", credsScript, script)
	}

	credsKey, ok := creds["script_key"]
	if !ok {
		t.Error("Creds map missing script_key")
	}
	if credsKey != key {
		t.Errorf("Script key miss-match: %v != %v", credsKey, key)
	}

}

func TestShotgunUser(t *testing.T) {
	host := "http://localhost"
	login := "fake-login"
	password := "fake-pass"
	sg := NewUserShotgun(
		host,
		login,
		password,
	)

	expectedUrl := host + "/api3/json"
	if sg.ServerUrl.String() != expectedUrl {
		t.Errorf("%v != %s", sg.ServerUrl.String())
	}

	creds := sg.Creds()

	credsLogin, ok := creds["user_login"]
	if !ok {
		t.Error("Creds map missing user_login")
	}
	if credsLogin != login {
		t.Errorf("User login miss-match: %v != %v", credsLogin, login)
	}

	credsPass, ok := creds["user_password"]
	if !ok {
		t.Error("Creds map missing user_password")
	}
	if credsPass != password {
		t.Errorf("User password miss-match: %v != %v", credsPass, password)
	}

}
