package main

import (
	"net/http"
	"proxy-forward/config"
	"proxy-forward/internal/http_proxy"
	"proxy-forward/internal/models"
	"proxy-forward/pkg/gcelery"
	"proxy-forward/pkg/gredis"
	"proxy-forward/pkg/logging"
)

func init() {
	models.Setup()
	gredis.Setup()
	gcelery.Setup()
}

func main() {
	goproxy := http_proxy.NewHandlerServer()
	logging.Log.Infof("Start the http server in port:%s", config.RuntimeViper.GetString("http_proxy_server.http_port"))
	go http.ListenAndServe(config.RuntimeViper.GetString("http_proxy_server.http_port"), goproxy.HttpHandler)
	logging.Log.Infof("Start the proxy server in port:%s", config.RuntimeViper.GetString("http_proxy_server.port"))
	logging.Log.Fatal(goproxy.ProxyHandler.ListenAndServe())
}
