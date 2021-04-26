package user_token_service

import (
	"encoding/json"
	"proxy-forward/internal/models"
	"proxy-forward/internal/service/cache_service"
	"proxy-forward/pkg/gredis"
	"proxy-forward/pkg/logging"
	"strconv"
)

type UserToken struct {
	ID              int
	Username        string
	Passwd          string
	ReqUsageAmount  int
	RespUsageAmount int
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
	gredis.Set(key, userToken, 120)
	return userToken, nil
}

func (u *UserToken) IncrReqBytes(length int) error {
	cache := cache_service.UserToken{ID: u.ID}
	key := cache.GetIncrReqKey()
	if gredis.Exists(key) {
		_, err := gredis.Incrby(key, length, 0)
		return err
	} else {
		_, err := gredis.Incrby(key, u.ReqUsageAmount, 0)
		_, err = gredis.Incrby(key, length, 0)
		return err
	}
}

func (u *UserToken) IncrRespBytes(length int) error {
	cache := cache_service.UserToken{ID: u.ID}
	key := cache.GetIncrRespKey()
	if gredis.Exists(key) {
		_, err := gredis.Incrby(key, length, 0)
		return err
	} else {
		_, err := gredis.Incrby(key, u.RespUsageAmount, 0)
		_, err = gredis.Incrby(key, length, 0)
		return err
	}
}

func (u *UserToken) SetReqUsageKey(id int) error {
	key := "REQ_USAGE_TOKEN_IDS"
	_, err := gredis.Sadd(key, strconv.Itoa(id))
	return err
}

func (u *UserToken) SetRespUsageKey(id int) error {
	key := "RESP_USAGE_TOKEN_IDS"
	_, err := gredis.Sadd(key, strconv.Itoa(id))
	return err
}
