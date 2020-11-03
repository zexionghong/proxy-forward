package cache_service

import (
	"proxy-forward/pkg/e"
	"strconv"
	"strings"
)

type IP struct {
	ID int
}

func (ip *IP) GetKey() string {
	keys := []string{
		e.CACHE_IP,
		strconv.Itoa(ip.ID),
	}
	return strings.Join(keys, "_")
}
