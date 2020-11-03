package user_token_service

import (
	"encoding/json"
	"proxy-forward/internal/models"
	"proxy-forward/internal/service/cache_service"
	"proxy-forward/pkg/gredis"
	"proxy-forward/pkg/logging"
)

type UserToken struct {
	Username string
	Passwd   string
}

func (u *UserToken) Get() (*models.UserToken, error) {
	var cacheUserToken *models.UserToken
	cache := cache_service.UserToken{Username: u.Username, Passwd: u.Passwd}
	key := cache.GetUserTokenKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Log.Warn(err.Error())
		} else {
			json.Unmarshal(data, &cacheUserToken)
			return cacheUserToken, err
		}
	}
	userToken, err := models.GetUserToken(u.Username, u.Passwd)
	if err != nil {
		return nil, err
	}
	gredis.Set(key, userToken, 60)
	return userToken, nil
}
