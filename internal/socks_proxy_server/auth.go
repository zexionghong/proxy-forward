package socks_proxy_server

import (
	"errors"
	"fmt"
	"proxy-forward/internal/models"
	"proxy-forward/internal/service/ip_service"
	"proxy-forward/internal/service/proxy_ip_service"
	"proxy-forward/internal/service/proxy_machine_service"
	"proxy-forward/internal/service/user_token_service"
	"proxy-forward/pkg/utils"
	"strconv"
)

// Auth provides basid authorization for handler server.
func Auth(username, password string) (*models.UserToken, bool) {
	userToken, ok := Verify(username, password)
	if !ok {
		return nil, false
	}
	return userToken, true
}

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
		return userToken, true
	}
	return nil, false
}

func LoadRemoteAddr(userToken *models.UserToken) (string, string, string, error) {
	var remoteAddr string
	if userToken == nil {
		return "", "", "", errors.New("load remote addr fail.")
	}
	if userToken.Expired == 1 {
		return "", "", "", errors.New("user token expired.")
	}
	proxyIPService := proxy_ip_service.ProxyIP{ID: userToken.PiID}
	proxyIP, err := proxyIPService.GetByID()
	if err != nil {
		return "", "", "", errors.New("load remote addr fail.")
	}
	if proxyIP.Online != 1 || proxyIP.Health != 1 || proxyIP.Status != 1 {
		return "", "", "", errors.New("load remote addr fail.")
	}
	proxyMachineService := proxy_machine_service.ProxyMachine{ID: proxyIP.PmID}
	proxyMachine, err := proxyMachineService.GetByID()
	if err != nil {
		return "", "", "", errors.New("load remote addr fail.")
	}
	ipService := ip_service.IP{ID: proxyMachine.IpID}
	iP, err := ipService.GetByID()
	if err != nil {
		return "", "", "", errors.New("load remote addr fail.")
	}
	if iP.IpAddress == "" {
		remoteAddr = utils.InetNtoA(iP.IpAddr)
	} else {
		ipAddress, err := strconv.ParseInt(iP.IpAddress, 16, 0)
		if err != nil {
			return "", "", "", errors.New("load remote addr fail.")
		}
		remoteAddr = utils.InetNtoA(ipAddress)
	}
	port := proxyIP.ForwardPort
	return fmt.Sprintf("%s:%d", remoteAddr, port), proxyIP.Username, proxyIP.Password, nil
}
