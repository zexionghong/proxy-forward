package handler

import "net/http"

func init() {
}

func (ps *ProxyServer) Done(rw http.ResponseWriter, req *http.Request) {
	// TODO: 流量统计
}
