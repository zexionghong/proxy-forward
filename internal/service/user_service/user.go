package user_service

import (
	"encoding/json"
	"proxy-forward/internal/models"
	"proxy-forward/internal/service/cache_service"
	"proxy-forward/pkg/gredis"
	"proxy-forward/pkg/logging"
)

type User struct {
	Uid int
}

func (u *User) Get() (*models.User, error) {
	var cacheUser *models.User
	cache := cache_service.User{Uid: u.Uid}
	key := cache.GetUserKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Log.Warn(err.Error())
		} else {
			json.Unmarshal(data, &cacheUser)
			return cacheUser, nil
		}
	}
	user, err := models.GetUser(u.Uid)
	if err != nil {
		return nil, err
	}
	gredis.Set(key, user, 120)
	return user, nil
}
