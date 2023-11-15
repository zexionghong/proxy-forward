package proxy

import (
	"fmt"
	"net/http"
	"net/url"
)

type ProxyServer struct {
	Travel   *http.Transport
	OnlyHttp bool
}

func NewProxyServer(remoteAddr string, port int, username, password string) (*ProxyServer, error) {

	if username == "" && password == "" {
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
	} else {
		u, err := url.Parse(fmt.Sprintf("http://%s:%s@%s:%d", username, password, remoteAddr, port))
		// u, err := url.Parse(fmt.Sprintf("http://%s:%d", remoteAddr, port))
		if err != nil {
			return nil, err
		}
		return &ProxyServer{
			OnlyHttp: true,
			Travel: &http.Transport{
				Proxy: http.ProxyURL(u),
				// TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}, nil

	}
}
