package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/Sirupsen/logrus"
	log "github.com/Sirupsen/logrus"
)

type Shotgun struct {
	ServerUrl    *url.URL
	ScriptName   string
	ScriptKey    string
	UserLogin    string
	UserPassword string
	client       http.Client
}

func NewShotgun(host, scriptName, scriptKey string) Shotgun {
	u, err := url.Parse(host)
	if err != nil {
		log.Fatalln(err)
	}
	u.Path = "/api3/json"

	return Shotgun{
		ServerUrl:  u,
		ScriptName: scriptName,
		ScriptKey:  scriptKey,
		client:     http.Client{},
	}
}

func NewUserShotgun(host, login, password string) Shotgun {
	u, err := url.Parse(host)
	if err != nil {
		log.Fatalln(err)
	}
	u.Path = "/api3/json"

	return Shotgun{
		ServerUrl:    u,
		UserLogin:    login,
		UserPassword: password,
		client:       http.Client{},
	}
}

func (sg *Shotgun) Creds() map[string]string {
	if sg.UserLogin != "" {
		log.WithFields(logrus.Fields{
			"login": sg.UserLogin,
		}).Debug("Using user credintials")
		return map[string]string{
			"user_login":    sg.UserLogin,
			"user_password": sg.UserPassword,
		}
	}
	log.WithFields(logrus.Fields{
		"login": sg.ScriptName,
	}).Debug("Using script credintials")
	return map[string]string{
		"script_name": sg.ScriptName,
		"script_key":  sg.ScriptKey,
	}
}

func (sg *Shotgun) Request(method_name string, query interface{}) (*http.Response, error) {
	requestData := make(map[string]interface{})
	requestData["method_name"] = method_name
	requestData["params"] = []interface{}{sg.Creds(), query}

	// requestData["params"][0] = sg.creds()
	// requestData["params"][1] =

	bodyJson, err := json.Marshal(requestData)
	log.Debug("Json Request:", string(bodyJson))
	if err != nil {
		log.Error("Shotugn.Request Marshal Error: ", err)
		return &http.Response{}, err
	}

	req, err := http.NewRequest("POST", sg.ServerUrl.String(), bytes.NewReader(bodyJson))
	if err != nil {
		log.Error("Shotugn.Request Http Request Error: ", err)
		return &http.Response{}, err
	}
	return sg.client.Do(req)
}
