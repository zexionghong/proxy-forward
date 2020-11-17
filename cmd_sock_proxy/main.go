package main

import (
	"os"
	"proxy-forward/config"
	"proxy-forward/internal/sock_proxy"
	"proxy-forward/pkg/logging"
)

func init() {
	// models.Setup()
	// gredis.Setup()
}

func main() {
	server := sock_proxy.NewServer()
	logging.Log.Infof("Start the socks5 proxy server in port:%s", config.RuntimeViper.GetString("sock_proxy_server.port"))
	if err := server.Run(); err != nil {
		logging.Log.Errorf("Run socks5 server error: %s", err.Error())
		os.Exit(1)
	}

	logging.Log.Info("Socks5 server normal exit.")
	os.Exit(0)
}
