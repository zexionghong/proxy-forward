package proxy

import (
	"net/http"

	px "golang.org/x/net/proxy"
)

type ProxyServer struct {
	Geo    string
	Travel *http.Transport
}

func NewProxyServer(geo, remoteAddr string) (*ProxyServer, error) {
	dialer, err := px.SOCKS5("tcp", remoteAddr, nil, px.Direct)
	if err != nil {
		return nil, err
	}
	return &ProxyServer{
		Geo:    geo,
		Travel: &http.Transport{Dial: dialer.Dial},
	}, nil

}
