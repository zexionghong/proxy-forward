package http_proxy

import (
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"proxy-forward/config"
	"proxy-forward/internal/http_proxy/proxy"
	"proxy-forward/internal/models"
	"proxy-forward/internal/service/user_token_service"
	"proxy-forward/pkg/logging"
	"time"

	"github.com/gin-gonic/gin"
)

type HttpProxyServer struct {
	ProxyHandler *http.Server
	HttpHandler  *gin.Engine
}

type HandlerServer struct {
}

// NewHandlerServer returns a new handler server.
func NewHandlerServer() *HttpProxyServer {

	return &HttpProxyServer{
		HttpHandler: InitRouter(),
		ProxyHandler: &http.Server{
			Addr:           config.RuntimeViper.GetString("http_proxy_server.port"),
			Handler:        &HandlerServer{},
			ReadTimeout:    time.Duration(config.RuntimeViper.GetInt("http_proxy_server.http_read_timeout")) * time.Second,
			WriteTimeout:   time.Duration(config.RuntimeViper.GetInt("http_proxy_server.http_write_timeout")) * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}
}

//ServeHTTP will be automatically called by system.
//HandlerServer implements the Handler interface which need ServeHTTP.
func (hs *HandlerServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logging.Log.Warnf("HandlerServer.ServeHTTP panic: %+v", err)
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}()

	userToken, ok := hs.Auth(rw, req)
	if !ok {
		return
	}
	// preload url 预加载对黑名单域名的检查
	ok = hs.PreloadReq(userToken, rw, req)
	if !ok {
		return
	}

	// load travel
	travel, ok := hs.LoadTraveling(userToken, rw, req)
	if !ok {
		return
	}
	defer hs.Done(rw, req)

	// http
	if travel.OnlyHttp == true {
		if req.Method == "CONNECT" {
			hs.OnlyHttpsHandler(travel, rw, req, userToken)
		} else {
			hs.OnlyHttpHandler(travel, rw, req, userToken)
		}
	} else { // sock
		if req.Method == "CONNECT" {
			hs.HttpsHandler(travel, rw, req, userToken)
		} else {
			hs.HttpHandler(travel, rw, req, userToken)
		}
	}
}

//HttpHandler handles http connections.
func (hs *HandlerServer) HttpHandler(travel *proxy.ProxyServer, rw http.ResponseWriter, req *http.Request, userToken *models.UserToken) {
	// sock request 字节数
	userTokenService := user_token_service.UserToken{ID: userToken.ID, ReqUsageAmount: userToken.ReqUsageAmount, RespUsageAmount: userToken.RespUsageAmount}
	reqBytes, _ := httputil.DumpRequest(req, true)
	//
	RmProxyHeaders(req)
	resp, err := travel.Travel.RoundTrip(req)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	defer resp.Body.Close()
	userTokenService.IncrReqBytes(len(reqBytes))
	userTokenService.SetReqUsageKey(userToken.ID)

	ClearHeaders(rw.Header())
	CopyHeaders(rw.Header(), resp.Header)

	rw.WriteHeader(resp.StatusCode)

	respBytes, _ := httputil.DumpResponse(resp, true)
	// sock reesponse 字节数
	userTokenService.IncrRespBytes(len(respBytes))
	userTokenService.SetRespUsageKey(userToken.ID)

	_, err = io.Copy(rw, resp.Body)
	if err != nil && err != io.EOF {
		return
	}
}

// HttpsHandler handles any connection which needs "connect" method.
func (hs *HandlerServer) HttpsHandler(travel *proxy.ProxyServer, rw http.ResponseWriter, req *http.Request, userToken *models.UserToken) {
	// sock request 字节数
	//
	hj, _ := rw.(http.Hijacker)
	Client, _, err := hj.Hijack()
	if err != nil {
		http.Error(rw, "Failed", http.StatusBadRequest)
		return
	}
	Remote, err := travel.Travel.Dial("tcp", req.URL.Host)
	if err != nil {
		http.Error(nil, "Failed", http.StatusBadGateway)
		return
	}
	Client.SetDeadline(time.Now().Add(time.Second * 60))
	Remote.SetDeadline(time.Now().Add(time.Second * 60))

	_, _ = Client.Write(HTTP200)
	go copyRemoteToClient(Remote, Client, userToken, 1)
	go copyRemoteToClient(Client, Remote, userToken, 2)
}

func copyRemoteToClient(Remote, Client net.Conn, userToken *models.UserToken, action int) {
	// aciont = 1  client => remote request
	// action = 2 remote => client response
	defer func() {
		_ = Remote.Close()
		_ = Client.Close()
	}()
	// 字节数
	userTokenService := user_token_service.UserToken{ID: userToken.ID, ReqUsageAmount: userToken.ReqUsageAmount, RespUsageAmount: userToken.RespUsageAmount}
	n, err := io.Copy(Remote, Client)
	if n > 0 {
		if action == 1 {
			userTokenService.IncrReqBytes(int(n))
			userTokenService.SetReqUsageKey(userToken.ID)
		} else if action == 2 {
			userTokenService.IncrRespBytes(int(n))
			userTokenService.SetRespUsageKey(userToken.ID)
		}
	}
	if err != nil && err != io.EOF {
		return
	}
}

// OnlyHttp proxy handles http connections
func (hs *HandlerServer) OnlyHttpHandler(travel *proxy.ProxyServer, rw http.ResponseWriter, req *http.Request, userToken *models.UserToken) {
	// request 字节数
	userTokenService := user_token_service.UserToken{ID: userToken.ID, ReqUsageAmount: userToken.ReqUsageAmount, RespUsageAmount: userToken.RespUsageAmount}
	reqBytes, _ := httputil.DumpRequest(req, true)
	//
	RmProxyHeaders(req)
	resp, err := travel.Travel.RoundTrip(req)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	defer resp.Body.Close()
	userTokenService.IncrReqBytes(len(reqBytes))
	userTokenService.SetReqUsageKey(userToken.ID)

	ClearHeaders(rw.Header())
	CopyHeaders(rw.Header(), resp.Header)

	rw.WriteHeader(resp.StatusCode)

	respBytes, _ := httputil.DumpResponse(resp, true)
	// response 字节数
	userTokenService.IncrRespBytes(len(respBytes))
	userTokenService.SetRespUsageKey(userToken.ID)
	_, err = io.Copy(rw, resp.Body)
	if err != nil && err != io.EOF {
		return
	}
}

// // OnlyHttpsHandler handlers any connection which needs "connect" method.
func (hs *HandlerServer) OnlyHttpsHandler(travel *proxy.ProxyServer, rw http.ResponseWriter, req *http.Request, userToken *models.UserToken) {
	RmProxyHeaders(req)
	hj, _ := rw.(http.Hijacker)
	Client, _, err := hj.Hijack()
	if err != nil {
		http.Error(rw, "Failed", http.StatusBadRequest)
		return
	}
	parnetUrl, err := travel.Travel.Proxy(req)
	if err != nil {
		http.Error(rw, "Failed", http.StatusBadRequest)
		return
	}
	Remote, err := net.Dial("tcp", parnetUrl.Host)
	if err != nil {
		logging.Log.Info(Remote, parnetUrl.Host, req.Host, err)
		http.Error(rw, "Failed", http.StatusBadGateway)
		return
	}
	username := parnetUrl.User.Username()
	password, _ := parnetUrl.User.Password()

	if username != "" && password != "" {
		auth := fmt.Sprintf("%s:%s", username, password)
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		_, _ = Remote.Write([]byte(fmt.Sprintf("CONNECT %s HTTP/1.1\r\nHost: %s\r\nProxy-Connection: Keey-Alive\r\nProxy-Authorization: %s\r\n\r\n", req.Host, req.Host, basicAuth)))
	} else {
		_, _ = Remote.Write([]byte(fmt.Sprintf("CONNECT %s HTTP/1.1\r\n\r\n", req.Host)))
	}

	go copyRemoteToClient(Remote, Client, userToken, 1)
	go copyRemoteToClient(Client, Remote, userToken, 2)
}
