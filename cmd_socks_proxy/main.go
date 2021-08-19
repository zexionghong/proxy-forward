package main

import (
	"os"
	"proxy-forward/config"
	"proxy-forward/internal/models"
	"proxy-forward/internal/socks_proxy_server"
	"proxy-forward/pkg/geolite"
	"proxy-forward/pkg/gredis"
	"proxy-forward/pkg/logging"
)

func init() {
	models.Setup()
	gredis.Setup()
	geolite.Setup()
}

func main() {
	server := socks_proxy_server.NewServer()
	logging.Log.Infof("Start the socks5 proxy server in port:%s", config.RuntimeViper.GetString("socks_proxy_server.port"))
	if err := server.Run(); err != nil {
		logging.Log.Errorf("Run socks5 server error: %s", err.Error())
		os.Exit(1)
	}

	logging.Log.Info("Socks5 server normal exit.")
	os.Exit(0)
}
