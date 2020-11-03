package cache_service

import (
	"proxy-forward/pkg/e"
	"strconv"
	"strings"
)

type ProxyIP struct {
	ID int
}

func (p *ProxyIP) GetKey() string {
	keys := []string{
		e.CACHE_PROXY_IP,
		strconv.Itoa(p.ID),
	}
	return strings.Join(keys, "_")
}
