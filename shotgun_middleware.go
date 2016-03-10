package main

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/context"
)

var connectionCache map[string]Shotgun

func init() {
	connectionCache = make(map[string]Shotgun)
}

func ShotgunAuthMiddleware(config clientConfig) func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	return func(rw http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		s := strings.SplitN(req.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			rw.Header().Set("WWW-Authenticate", `Basic realm="shotgun-restful"`)
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("401 Unauthorized\n"))
			return
		}

		if !strings.HasPrefix(s[0], "Basic") {
			rw.Header().Set("WWW-Authenticate", `Basic realm="shotgun-restful"`)
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("401 Unauthorized\n"))
			return
		}

		var isUser bool
		if s[0] == "Basic-User" {
			isUser = true
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			rw.WriteHeader(http.StatusForbidden)
			return
		}

		name := pair[0]
		key := pair[1]

		hasher := sha1.New()
		hash := hex.EncodeToString(hasher.Sum([]byte(fmt.Sprintf("%s%v%s%s", config.shotgunHost, isUser, name, key))))

		conn, ok := connectionCache[hash]
		if !ok {
			if isUser {
				conn = NewUserShotgun(config.shotgunHost, name, key)
			} else {
				conn = NewShotgun(config.shotgunHost, name, key)
			}

		}
		connectionCache[hash] = conn
		conn.Log()
		context.Set(req, "sgConn", conn)

		next(rw, req)
	}
}
