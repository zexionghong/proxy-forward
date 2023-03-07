package user_token_service

import (
	"encoding/json"
	"proxy-forward/internal/models"
	"proxy-forward/internal/service/cache_service"

	// "proxy-forward/pkg/gcelery"
	"proxy-forward/pkg/gmongo"
	"proxy-forward/pkg/gredis"
	"proxy-forward/pkg/logging"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type UserToken struct {
	ID              int
	Username        string
	Passwd          string
	ReqUsageAmount  int
	RespUsageAmount int
	PsID            int
	IsStatic        int
	DataCenter      int
	LaID            int
	Uid             int
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

func (u *UserToken) IncrReqBytes(remoteAddr string, length int) error {
	// usage statistics
	// data := map[string]interface{}{
	data := bson.M{
		"remote_addr":   remoteAddr,
		"usage":         length,
		"user_token_id": u.ID,
		"type":          "req",
		"ps_id":         u.PsID,
		"la_id":         u.LaID,
		"is_static":     u.IsStatic,
		"data_center":   u.DataCenter,
		"uid":           u.Uid,
		"timestamp":     int64(time.Now().Unix()),
	}
	gmongo.SaveForwardData(data)
	// gcelery.SendForwardDataTask(data)
	cache := cache_service.UserToken{ID: u.ID}
	key := cache.GetIncrReqKey()
	if gredis.Exists(key) {
		_, err := gredis.Incrby(key, length, 0)
		return err
	} else {
		lock_key := "LOCK_" + key
		lock, err := gredis.Incr(lock_key, 3)
		if err != nil {
			return err
		}
		if lock == 1 {
			gredis.Delete(key)
			_, err := gredis.Incrby(key, u.ReqUsageAmount+length, 0)
			gredis.Expired(lock_key, 1)
			return err
		}
		return nil
	}
}

func (u *UserToken) IncrRespBytes(remoteAddr string, length int) error {
	// usage statistics
	// data := map[string]interface{}{
	data := bson.M{
		"remote_addr":   remoteAddr,
		"usage":         length,
		"user_token_id": u.ID,
		"type":          "resp",
		"ps_id":         u.PsID,
		"la_id":         u.LaID,
		"is_static":     u.IsStatic,
		"data_center":   u.DataCenter,
		"uid":           u.Uid,
		"timestamp":     int64(time.Now().Unix()),
	}
	gmongo.SaveForwardData(data)
	// gcelery.SendForwardDataTask(data)
	cache := cache_service.UserToken{ID: u.ID}
	key := cache.GetIncrRespKey()
	if gredis.Exists(key) {
		_, err := gredis.Incrby(key, length, 0)
		return err
	} else {
		lock_key := "LOCK_" + key
		lock, err := gredis.Incr(lock_key, 3)
		if err != nil {
			return err
		}
		if lock == 1 {
			gredis.Delete(key)
			_, err := gredis.Incrby(key, u.RespUsageAmount+length, 0)
			gredis.Expired(lock_key, 1)
			return err
		}
		return nil
	}
}

func (u *UserToken) SetReqUsageKey(id int) error {
	// key := "REQ_USAGE_TOKEN_IDS"
	// _, err := gredis.Sadd(key, strconv.Itoa(id))
	key := "REQ_USAGE_TOKEN_SORTED_IDS"
	_, err := gredis.ZaddByInt(key, int(time.Now().Unix()), strconv.Itoa(id))
	return err
}

func (u *UserToken) SetRespUsageKey(id int) error {
	// key := "RESP_USAGE_TOKEN_IDS"
	// _, err := gredis.Sadd(key, strconv.Itoa(id))
	key := "RESP_USAGE_TOKEN_SORTED_IDS"
	_, err := gredis.ZaddByInt(key, int(time.Now().Unix()), strconv.Itoa(id))
	return err
}
