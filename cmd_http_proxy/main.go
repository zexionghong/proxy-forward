package main

import (
	"proxy-forward/config"
	"proxy-forward/internal/http_proxy"
	"proxy-forward/internal/models"
	"proxy-forward/pkg/gredis"
	"proxy-forward/pkg/logging"
)

func init() {
	models.Setup()
	_ = gredis.Setup()
	//gcelery.Setup()
	//gmongo.Setup()
}

func main() {
	goproxy := http_proxy.NewHandlerServer()
	logging.Log.Errorf("Start the http server in port:%s", config.RuntimeViper.GetString("http_proxy_server.http_port"))
	logging.Log.Infof("Start the proxy server in port:%s", config.RuntimeViper.GetString("http_proxy_server.port"))
	//go http.ListenAndServe(config.RuntimeViper.GetString("http_proxy_server.http_port"), goproxy.HttpHandler)
	logging.Log.Fatal(goproxy.ProxyHandler.ListenAndServe())
}
