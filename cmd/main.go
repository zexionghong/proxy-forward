package main

import (
	"proxy-forward/config"
	"proxy-forward/internal/handler"
	"proxy-forward/pkg/logging"
)

func init() {
	logging.Setup()
	// models.Setup()
}

func main() {
	goproxy := handler.NewProxyServer()
	logging.Log.Infof("Start the proxy server in port:%s", config.RuntimeViper.GetString("server.port"))
	logging.Log.Fatal(goproxy.ListenAndServe())
}
