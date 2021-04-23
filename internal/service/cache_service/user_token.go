package cache_service

import (
	"proxy-forward/pkg/e"
	"strconv"
	"strings"
)

type UserToken struct {
	ID       int
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

func (u *UserToken) GetIncrReqKey() string {
	keys := []string{
		e.CACHE_USER_TOKEN_REQ_BYTES,
		strconv.Itoa(u.ID),
	}
	return strings.Join(keys, "_")
}

func (u *UserToken) GetIncrRespKey() string {
	keys := []string{
		e.CACHE_USER_TOKEN_RESP_BYTES,
		strconv.Itoa(u.ID),
	}
	return strings.Join(keys, "_")
}
