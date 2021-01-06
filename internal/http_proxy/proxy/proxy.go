package proxy

import (
	"fmt"
	"net/http"
	"net/url"

	px "golang.org/x/net/proxy"
)

type ProxyServer struct {
	Travel   *http.Transport
	OnlyHttp bool
}

func NewProxyServer(remoteAddr string, port int, username, password string, onlyHttp int) (*ProxyServer, error) {
	var (
		dialer px.Dialer
		err    error
	)
	if onlyHttp == 0 {
		if username == "" && password == "" {
			dialer, err = px.SOCKS5("tcp", fmt.Sprintf("%s:%d", remoteAddr, port), nil, px.Direct)
			if err != nil {
				return nil, err
			}
		} else {
			dialer, err = px.SOCKS5("tcp", fmt.Sprintf("%s:%d", remoteAddr, port), &px.Auth{username, password}, px.Direct)
			if err != nil {
				return nil, err
			}
		}
		return &ProxyServer{
			OnlyHttp: false,
			Travel:   &http.Transport{Dial: dialer.Dial},
		}, nil
	} else {
		u, err := url.Parse(fmt.Sprintf("http://%s:%d", remoteAddr, port))
		if err != nil {
			return nil, err
		}
		return &ProxyServer{
			OnlyHttp: true,
			Travel: &http.Transport{
				Proxy: http.ProxyURL(u),
			},
		}, nil
	}
}
