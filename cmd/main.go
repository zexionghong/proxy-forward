package main

import (
	"proxy-forward/config"
	"proxy-forward/internal/handler"
	"proxy-forward/internal/models"
	"proxy-forward/pkg/logging"
)

func init() {
	models.Setup()
}

func main() {
	goproxy := handler.NewHandlerServer()
	logging.Log.Infof("Start the proxy server in port:%s", config.RuntimeViper.GetString("server.port"))
	logging.Log.Fatal(goproxy.ListenAndServe())
}
