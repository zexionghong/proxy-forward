package http_proxy

import (
	"fmt"
	"net/http"
	"proxy-forward/internal/http_proxy/proxy"
	"proxy-forward/internal/models"
	"proxy-forward/internal/service/ip_service"
	"proxy-forward/internal/service/proxy_ip_service"
	"proxy-forward/internal/service/proxy_machine_service"
	"proxy-forward/pkg/logging"
	"proxy-forward/pkg/utils"
)

func init() {
}

const (
	CACHE_CAMP_PROXY = "CAMP_PROXY"
)

// Load username:password match ip:port sock connection
func (hs *HandlerServer) LoadTraveling(userToken *models.UserToken, rw http.ResponseWriter, req *http.Request) (*proxy.ProxyServer, bool) {
	travel, err := hs.loadTraveling(userToken, rw, req)
	if err != nil {
		return nil, false
	}
	return travel, true
}

// Load username:password match ip:port sock connection
func (hs *HandlerServer) loadTraveling(userToken *models.UserToken, rw http.ResponseWriter, req *http.Request) (*proxy.ProxyServer, error) {
	if userToken == nil {
		Unavailable(rw)
	}
	_cacheKey := fmt.Sprintf("%s_%d", CACHE_CAMP_PROXY, userToken.PiID)
	if tmp, ok := Camp.Get(_cacheKey); ok {
		travel := tmp.(*proxy.ProxyServer)
		return travel, nil
	}

	var (
		remoteAddr string
		port       int
	)
	proxyIPService := proxy_ip_service.ProxyIP{ID: userToken.PiID}
	proxyIP, err := proxyIPService.GetByID()
	if err != nil {
		Unavailable(rw)
		return nil, err
	}
	if proxyIP.Online != 1 && proxyIP.Health != 1 {
		Unavailable(rw)
		return nil, err
	}
	proxyMachineService := proxy_machine_service.ProxyMachine{ID: proxyIP.PmID}
	proxyMachine, err := proxyMachineService.GetByID()
	if err != nil {
		Unavailable(rw)
		return nil, err
	}
	ipService := ip_service.IP{ID: proxyMachine.IpID}
	iP, err := ipService.GetByID()
	if err != nil {
		Unavailable(rw)
		return nil, err
	}
	remoteAddr = utils.InetNtoA(iP.IpAddr)
	port = proxyIP.ForwardPort
	travel, ok := Connection(remoteAddr, port)
	if !ok {
		Unavailable(rw)
		return nil, err
	}
	Camp.Set(_cacheKey, travel)
	return travel, nil
}

// TODO: 流量统计 请求统计
func (hs *HandlerServer) Done(rw http.ResponseWriter, req *http.Request) {
}

func Unavailable(rw http.ResponseWriter) {
	hj, _ := rw.(http.Hijacker)
	Client, _, err := hj.Hijack()
	defer Client.Close()
	if err != nil {
		logging.Log.Warnf("fail to get TCP connection of client in Unavailable, %v", err)
	}
	_, _ = Client.Write(HTTP503)
}

// build tcp connection to remoteAddr:port
func Connection(remoteAddr string, port int) (*proxy.ProxyServer, bool) {
	if remoteAddr == "" || port == 0 {
		return nil, false
	}
	proxyServer, err := proxy.NewProxyServer(remoteAddr, port)
	if err != nil {
		return nil, false
	}
	return proxyServer, true
}
