package http_proxy

import (
	"encoding/base64"
	"errors"
	"net/http"
	"proxy-forward/internal/models"
	"proxy-forward/internal/service/user_service"
	"proxy-forward/internal/service/user_token_service"
	"proxy-forward/pkg/logging"
	"strings"
)

var (
	ERR_PROXY_AUTH = errors.New("fail to decoding Proxy-Authorization")
	ERR_LOGIN_IN   = errors.New("fail to login")
)

// Auth provides basic authorization for handler server.
func (hs *HandlerServer) Auth(rw http.ResponseWriter, req *http.Request) (*models.UserToken, bool) {
	userToken, err := hs.auth(rw, req)
	if err != nil {
		return nil, false
	}
	return userToken, true
}

// Auth provides basic authorization for handler server.
func (hs *HandlerServer) auth(rw http.ResponseWriter, req *http.Request) (*models.UserToken, error) {
	auth := req.Header.Get("Proxy-Authorization")
	auth = strings.Replace(auth, "Basic ", "", 1)
	if auth == "" {
		NeedAuth(rw)
	}
	data, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return nil, ERR_PROXY_AUTH
	}

	var username, password string

	userPasswordPair := strings.Split(string(data), ":")
	if len(userPasswordPair) != 2 {
		return nil, ERR_LOGIN_IN
	}
	username = userPasswordPair[0]
	var suffix string
	if strings.Contains(username, "-") {
		temp := strings.Split(username, "-")
		username = temp[0]
		if len(temp) > 1 {
			suffix = "-" + strings.Join(temp[1:], "-")
		}
	}
	password = userPasswordPair[1]
	userToken, ok := Verify(username, password)
	if !ok {
		NeedAuth(rw)
		return nil, ERR_LOGIN_IN
	}
	userToken.Suffix = suffix
	return userToken, nil
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
func Verify(username, password string) (*models.UserToken, bool) {
	if username != "" && password != "" {
		userTokenService := user_token_service.UserToken{Username: username, Passwd: password}
		userToken, err := userTokenService.Get()
		if err != nil {
			return nil, false
		}
		if userToken.ID == 0 {
			return nil, false
		}
		userService := user_service.User{Uid: userToken.Uid}
		user, err := userService.Get()
		if err != nil {
			return nil, false
		}
		if !user.IsUse {
			return nil, false
		}
		return userToken, true
	}
	return nil, false
}
