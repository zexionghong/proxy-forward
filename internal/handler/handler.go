package handler

import (
	"io"
	"net"
	"net/http"
	"proxy-forward/config"
	"time"

	cmap "github.com/orcaman/concurrent-map"
)

type HandlerServer struct {
	//
	Camp cmap.ConcurrentMap
}

// NewHandlerServer returns a new handler server.
func NewHandlerServer() *http.Server {
	return &http.Server{
		Addr:           config.RuntimeViper.GetString("server.port"),
		Handler:        &HandlerServer{Camp: cmap.New()},
		ReadTimeout:    time.Duration(config.RuntimeViper.GetInt("server.http_read_timeout")) * time.Second,
		WriteTimeout:   time.Duration(config.RuntimeViper.GetInt("server.http_write_timeout")) * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}

//ServeHTTP will be automatically called by system.
//HandlerServer implements the Handler interface which need ServeHTTP.
func (hs *HandlerServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}()

	// TODO http Authorization
	if !hs.Auth(rw, req) {
		return
	}

	// load travel
	defer hs.Done(rw, req)

	if req.Method == "CONNECT" {
		// ps.HttpsHandler(rw, req)
	} else {
		// ps.HttpHandler(rw, req)
	}
}

//HttpHandler handles http connections.
func (hs *HandlerServer) HttpHandler(travel *http.Transport, rw http.ResponseWriter, req *http.Request) {
	RmProxyHeaders(req)
	resp, err := travel.RoundTrip(req)
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
func (hs *HandlerServer) HttpsHandler(travel *http.Transport, rw http.ResponseWriter, req *http.Request) {
	hj, _ := rw.(http.Hijacker)
	Client, _, err := hj.Hijack()
	if err != nil {
		http.Error(rw, "Failed", http.StatusBadRequest)
		return
	}
	Remote, err := travel.Dial("tcp", req.URL.Host)
	if err != nil {
		http.Error(rw, "Failed", http.StatusBadGateway)
		return
	}

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
