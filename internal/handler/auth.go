package handler

import (
	"encoding/base64"
	"errors"
	"net/http"
	"proxy-forward/config"
	"proxy-forward/pkg/logging"
	"strings"
)

var (
	ERR_PROXY_AUTH = errors.New("fail to decoding Proxy-Authorization")
	ERR_LOGIN_IN   = errors.New("fail to login")
)

//Auth provides basic authorization for proxy server.
func (ps *ProxyServer) Auth(rw http.ResponseWriter, req *http.Request) bool {
	var err error
	if config.RuntimeViper.GetBool("server.auth") {
		if err = ps.auth(rw, req); err != nil {
			return false
		}
		return true
	}
	return true
}

// Auth provides basic authorization for proxy server.
func (ps *ProxyServer) auth(rw http.ResponseWriter, req *http.Request) (string, error) {
	auth := req.Header.Get("Proxy-Authorization")
	auth = strings.Replace(auth, "Basic ", "", 1)
	if auth == "" {
		NeedAuth(rw)
	}
	data, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return "", ERR_PROXY_AUTH
	}

	var user, password string

	userPasswordPair := strings.Split(string(data), ":")
	if len(userPasswordPair) != 2 {
		return "", ERR_LOGIN_IN
	}
	user = userPasswordPair[0]
	password = userPasswordPair[1]
	if Verify(user, password) == false {
		NeedAuth(rw)
		return "", ERR_LOGIN_IN
	}
	return user, nil
}

func NeedAuth(rw http.ResponseWriter) {
	hj, _ := rw.(http.Hijacker)
	Client, _, err := hj.Hijack()
	defer Client.Close()
	if err != nil {
		logging.Log.Warnf("fail to get TCP connection of client in auth, %v", err)
	}
	_, _ = Client.Write(HTTP407)
}

// Verify verifies username and password
func Verify(user, password string) bool {
	if user != "" && password != "" {
		return true
	}
	return false
}
