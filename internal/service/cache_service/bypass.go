package cache_service

import (
	"proxy-forward/pkg/e"
	"strconv"
	"strings"
)

type Bypass struct {
	PsID int
}

func (b *Bypass) GetKeyByPsID() string {
	keys := []string{
		e.CACHE_BYPASS,
		"PSID",
		strconv.Itoa(b.PsID),
	}
	return strings.Join(keys, "_")
}
