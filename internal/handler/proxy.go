package handler

import (
	"net/http"
	"proxy-forward/config"
	"time"

	cmap "github.com/orcaman/concurrent-map"
)

type ProxyServer struct {
	//
	Camp cmap.ConcurrentMap
}

// NewProxyServer returns a new proxy server.
func NewProxyServer() *http.Server {
	return &http.Server{
		Addr:           config.RuntimeViper.GetString("server.port"),
		Handler:        &ProxyServer{Camp: cmap.New()},
		ReadTimeout:    time.Duration(config.RuntimeViper.GetInt("server.http_read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(config.RuntimeViper.GetInt("server.http_write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

//ServeHTTP will be automatically called by system.
//ProxyServer implements the Handler interface which need ServeHTTP.
func (ps *ProxyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}()

	// TODO http Authorization
	if !ps.Auth(rw, req) {
		return
	}

	if req.Method == "CONNECT" {
		ps.HttpsHandler(rw, req)
	} else {
		ps.HttpHandler(rw, req)
	}
}

//HttpHandler handles http connections.
func (ps *ProxyServer) HttpHandler(rw http.ResponseWriter, req *http.Request) {
	RmProxyHeaders(req)

}

// HttpsHandler handles any connection which needs "connect" method.
func (ps *ProxyServer) HttpsHandler(rw http.ResponseWriter, req *http.Request) {

}
