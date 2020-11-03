package cache_service

import (
	"proxy-forward/pkg/e"
	"strconv"
	"strings"
)

type ProxyMachine struct {
	ID int
}

func (p *ProxyMachine) GetKey() string {
	keys := []string{
		e.CACHE_PROXY_MACHINE,
		strconv.Itoa(p.ID),
	}
	return strings.Join(keys, "_")
}
