package cache_service

import (
	"proxy-forward/pkg/e"
	"strconv"
	"strings"
)

type User struct {
	Uid int
}

func (u *User) GetUserKey() string {
	keys := []string{
		e.CACHE_USER,
		strconv.Itoa(u.Uid),
	}
	return strings.Join(keys, "_")
}
