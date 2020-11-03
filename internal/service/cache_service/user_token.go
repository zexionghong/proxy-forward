package cache_service

import (
	"proxy-forward/pkg/e"
	"strings"
)

type UserToken struct {
	Username string
	Passwd   string
}

func (u *UserToken) GetUserTokenKey() string {
	keys := []string{
		e.CACHE_USER_TOKEN,
		u.Username,
		u.Passwd,
	}
	return strings.Join(keys, "_")
}
