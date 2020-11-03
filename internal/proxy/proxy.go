package proxy

import (
	"fmt"
	"net/http"

	px "golang.org/x/net/proxy"
)

type ProxyServer struct {
	Travel *http.Transport
}

func NewProxyServer(remoteAddr string, port int) (*ProxyServer, error) {
	dialer, err := px.SOCKS5("tcp", fmt.Sprintf("%s:%d", remoteAddr, port), nil, px.Direct)
	if err != nil {
		return nil, err
	}
	return &ProxyServer{
		Travel: &http.Transport{Dial: dialer.Dial},
	}, nil

}
