package cache_service

import (
	"proxy-forward/pkg/e"
	"strconv"
	"strings"
)

type BypassProxySupplier struct {
	PsID int
}

func (bps *BypassProxySupplier) GetKeyByPsID() string {
	keys := []string{
		e.CACHE_BYPASS_PROXY_SUPPLIER,
		strconv.Itoa(bps.PsID),
	}
	return strings.Join(keys, "_")
}
