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
	gredis.Setup()
}

func main() {
	goproxy := http_proxy.NewHandlerServer()
	logging.Log.Infof("Start the proxy server in port:%s", config.RuntimeViper.GetString("server.port"))
	logging.Log.Fatal(goproxy.ProxyHandler.ListenAndServe())
}
