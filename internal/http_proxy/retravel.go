package http_proxy

import (
	"errors"
	"net/http"
	"proxy-forward/internal/http_proxy/proxy"
	"proxy-forward/internal/models"
	"proxy-forward/internal/service/proxy_ip_service"
	"proxy-forward/pkg/logging"
)

func init() {
}

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
		return nil, errors.New("userToken is nil")
	}
	if userToken.Expired {
		Unavailable(rw)
		return nil, errors.New("userToken was expired")
	}
	if userToken.IsDeleted {
		Unavailable(rw)
		return nil, errors.New("proxy unavailable")
	}

	if !userToken.IsUse {
		Unavailable(rw)
		return nil, errors.New("token close")
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
	//if proxyIP.Online != 1 || proxyIP.Health != 1 || proxyIP.Status != 1 {
	//	Unavailable(rw)
	//	return nil, errors.New("proxy unavailable")
	//}
	//proxyMachineService := proxy_machine_service.ProxyMachine{ID: proxyIP.PmID}
	//proxyMachine, err := proxyMachineService.GetByID()
	//if err != nil {
	//	Unavailable(rw)
	//	return nil, err
	//}
	//proxySupplierService := proxy_supplier_service.ProxySupplier{ID: proxyMachine.PsID}
	//proxySupplier, err := proxySupplierService.GetByID()
	//if err != nil {
	//	Unavailable(rw)
	//	return nil, err
	//}
	//// luminati 直接调用域名
	//if proxyMachine.IpID == 0 {
	//	if proxyMachine.Domain == "" {
	//		Unavailable(rw)
	//		return nil, err
	//	}
	//	remoteAddr = proxyMachine.Domain
	//} else {
	//	ipService := ip_service.IP{ID: proxyMachine.IpID}
	//	iP, err := ipService.GetByID()
	//	if err != nil {
	//		Unavailable(rw)
	//		return nil, err
	//	}
	//	if iP.IpAddress == "" {
	//		remoteAddr = utils.InetNtoA(iP.IpAddr)
	//	} else {
	//		ipAddress, err := strconv.ParseInt(iP.IpAddress, 16, 0)
	//		if err != nil {
	//			return nil, err
	//		}
	//		remoteAddr = utils.InetNtoA(ipAddress)
	//	}
	//}
	remoteAddr = proxyIP.Host

	var username, password string

	username = proxyIP.Username
	password = proxyIP.Password

	port = proxyIP.Port
	travel, ok := Connection(remoteAddr, port, username, password)
	if !ok {
		Unavailable(rw)
		return nil, err
	}
	return travel, nil
}

// TODO: 流量统计 请求统计
func (hs *HandlerServer) Done(rw http.ResponseWriter, req *http.Request) {
	logging.Log.Infof("请求结束 URL:%s", req.URL)
}

func Unavailable(rw http.ResponseWriter) {
	hj, _ := rw.(http.Hijacker)
	Client, _, err := hj.Hijack()
	defer Client.Close()
	if err != nil {
		logging.Log.Warnf("fail to get TCP connection of client in Unavailable, %v", err)
	}
	_, _ = Client.Write(HTTP504)
}

// build tcp connection to remoteAddr:port
func Connection(remoteAddr string, port int, username, password string) (*proxy.ProxyServer, bool) {
	if remoteAddr == "" || port == 0 {
		return nil, false
	}
	proxyServer, err := proxy.NewProxyServer(remoteAddr, port, username, password)
	if err != nil {
		return nil, false
	}
	return proxyServer, true
}
