package http_proxy

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"proxy-forward/config"
	"proxy-forward/internal/http_proxy/proxy"
	"proxy-forward/pkg/logging"
	"time"

	"github.com/gin-gonic/gin"
	cmap "github.com/orcaman/concurrent-map"
)

type HttpProxyServer struct {
	ProxyHandler *http.Server
	HttpHandler  *gin.Engine
}

type HandlerServer struct {
}

var (
	Camp cmap.ConcurrentMap
)

// NewHandlerServer returns a new handler server.
func NewHandlerServer() *HttpProxyServer {
	Camp = cmap.New()

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

	// load travel
	travel, ok := hs.LoadTraveling(userToken, rw, req)
	if !ok {
		return
	}
	defer hs.Done(rw, req)

	// http
	if travel.OnlyHttp == true {
		if req.Method == "CONNECT" {
			hs.OnlyHttpsHandler(travel, rw, req)
		} else {
			hs.OnlyHttpHandler(travel, rw, req)
		}
	} else { // sock
		if req.Method == "CONNECT" {
			hs.HttpsHandler(travel, rw, req)
		} else {
			hs.HttpHandler(travel, rw, req)
		}
	}
}

//HttpHandler handles http connections.
func (hs *HandlerServer) HttpHandler(travel *proxy.ProxyServer, rw http.ResponseWriter, req *http.Request) {
	RmProxyHeaders(req)
	resp, err := travel.Travel.RoundTrip(req)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	ClearHeaders(rw.Header())
	CopyHeaders(rw.Header(), resp.Header)

	rw.WriteHeader(resp.StatusCode)

	_, err = io.Copy(rw, resp.Body)
	if err != nil && err != io.EOF {
		return
	}
}

// HttpsHandler handles any connection which needs "connect" method.
func (hs *HandlerServer) HttpsHandler(travel *proxy.ProxyServer, rw http.ResponseWriter, req *http.Request) {
	hj, _ := rw.(http.Hijacker)
	Client, _, err := hj.Hijack()
	if err != nil {
		http.Error(rw, "Failed", http.StatusBadRequest)
		return
	}
	Remote, err := travel.Travel.Dial("tcp", req.URL.Host)
	if err != nil {
		log.Println(Remote, err)
		http.Error(nil, "Failed", http.StatusBadGateway)
		return
	}
	Client.SetDeadline(time.Now().Add(time.Second * 10))
	Remote.SetDeadline(time.Now().Add(time.Second * 10))

	_, _ = Client.Write(HTTP200)
	go copyRemoteToClient(Remote, Client)
	go copyRemoteToClient(Client, Remote)
}

func copyRemoteToClient(Remote, Client net.Conn) {
	defer func() {
		_ = Remote.Close()
		_ = Client.Close()
	}()

	_, err := io.Copy(Remote, Client)
	if err != nil && err != io.EOF {
		return
	}
}

// OnlyHttp proxy handles http connections
func (hs *HandlerServer) OnlyHttpHandler(travel *proxy.ProxyServer, rw http.ResponseWriter, req *http.Request) {
	RmProxyHeaders(req)
	resp, err := travel.Travel.RoundTrip(req)
	if err != nil {
		http.Error(rw, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	ClearHeaders(rw.Header())
	CopyHeaders(rw.Header(), resp.Header)

	rw.WriteHeader(resp.StatusCode)

	_, err = io.Copy(rw, resp.Body)
	if err != nil && err != io.EOF {
		return
	}
}

// // OnlyHttpsHandler handlers any connection which needs "connect" method.
func (hs *HandlerServer) OnlyHttpsHandler(travel *proxy.ProxyServer, rw http.ResponseWriter, req *http.Request) {
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
		http.Error(rw, "Failed", http.StatusBadGateway)
		return
	}
	_, _ = Remote.Write([]byte(fmt.Sprintf("CONNECT %s HTTP/1.1\r\n\r\n", req.Host)))

	go copyRemoteToClient(Remote, Client)
	go copyRemoteToClient(Client, Remote)
}
