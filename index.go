package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/context"
)

func indexHandler(rw http.ResponseWriter, req *http.Request) {
	// For tesing more than anyting. It should hit the else 99% of the
	// time outside of testing.
	sg_conn, ok := context.GetOk(req, "sg_conn")
	var sg Shotgun
	if ok {
		sg = sg_conn.(Shotgun)
	} else {
		sg = NewShotgun(SG_HOST, "fake-script", "fake-key")
	}

	sgReq, err := sg.Request("info", make(map[string]interface{}))
	if err != nil {
		log.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if sgReq.StatusCode >= 400 && sgReq.StatusCode < 500 {
		log.Errorf("Shotgun Response Status Code: %v ", sgReq.StatusCode)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	} else if sgReq.StatusCode >= 500 {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	var infoResp map[string]interface{}
	respBody, err := ioutil.ReadAll(sgReq.Body)
	if err != nil {
		log.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(respBody, &infoResp)
	if err != nil {
		log.Error(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	versionSlice := infoResp["version"].([]interface{})

	infoJson := make(map[string]interface{})
	infoJson["shotgun_version"] = fmt.Sprintf("v%d.%d.%d",
		int(versionSlice[0].(float64)),
		int(versionSlice[1].(float64)),
		int(versionSlice[2].(float64)))

	infoJson["rest_version"] = VERSION

	encoder := json.NewEncoder(rw)
	err = encoder.Encode(infoJson)
	if err != nil {
		log.Error("Error encoding info json.")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

}
