package cache_service

import (
	"proxy-forward/pkg/e"
	"strconv"
	"strings"
)

type ProxySupplier struct {
	ID int
}

func (p *ProxySupplier) GetKey() string {
	keys := []string{
		e.CACHE_PROXY_SUPPLIER,
		strconv.Itoa(p.ID),
	}
	return strings.Join(keys, "_")
}
