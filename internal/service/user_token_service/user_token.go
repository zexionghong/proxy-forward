package user_token_service

import (
	"encoding/json"
	"fmt"
	"proxy-forward/internal/models"
	"proxy-forward/internal/service/cache_service"
	"proxy-forward/internal/service/proxy_ip_service"

	// "proxy-forward/pkg/gcelery"
	"proxy-forward/pkg/coarsetime"
	"proxy-forward/pkg/e"
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
	score, _ := gredis.Zscore(key, strconv.Itoa(id))
	if score > 0 {
		return nil
	} else {
		_, err := gredis.ZaddByInt(key, int(time.Now().Unix()), strconv.Itoa(id))
		return err
	}
}

func (u *UserToken) SetRespUsageKey(id int) error {
	// key := "RESP_USAGE_TOKEN_IDS"
	// _, err := gredis.Sadd(key, strconv.Itoa(id))
	key := "RESP_USAGE_TOKEN_SORTED_IDS"
	score, _ := gredis.Zscore(key, strconv.Itoa(id))
	if score > 0 {
		return nil
	} else {
		_, err := gredis.ZaddByInt(key, int(time.Now().Unix()), strconv.Itoa(id))
		return err
	}
}

func (u *UserToken) CollectReqUsage(id int, piID int, isStatic int, dataCenter int, length int) error {
	proxyIPService := proxy_ip_service.ProxyIP{ID: piID}
	proxyIP, err := proxyIPService.GetByID()
	if err != nil {
		return nil
	}
	var incrReqKey = ""
	today := coarsetime.CeilingTimezoneTimeNowYYMMDD(0)
	if isStatic == 0 {
		// 动态代理
		if proxyIP.ForwardPort == 22225 { // luminati 线路2
			incrReqKey = fmt.Sprintf("%s_%d", e.CACHE_LUMINATI_RESIDENTIAL_USAGE_REQ, today)
		} else if proxyIP.ForwardPort == 7777 { // oxylab 线路1
			incrReqKey = fmt.Sprintf("%s_%d", e.CACHE_OXYLAB_RESIDENTIAL_USAGE_REQ, today)
		} else {
			incrReqKey = fmt.Sprintf("%s_%d", e.CACHE_922_RESIDENTIAL_USAGE_REQ, today)
		}
	} else if isStatic == 1 {
		// 静态代理
		if dataCenter == 0 {
			if proxyIP.ForwardPort == 22225 { // luminati
				incrReqKey = fmt.Sprintf("%s_%d", e.CACHE_LUMINATI_ISP_USAGE_REQ, today)
			} else { // iproyal
				incrReqKey = fmt.Sprintf("%s_%d", e.CACHE_IPROYAL_ISP_USAGE_REQ, today)
			}
		} else if dataCenter == 1 { //机房代理
			if proxyIP.ForwardPort == 22225 {
				incrReqKey = fmt.Sprintf("%s_%d", e.CACHE_LUMINATI_DATACENTER_USAGE_REQ, today)
			} else { // instant_proxy
				incrReqKey = fmt.Sprintf("%s_%d", e.CACHE_INSTANTPROXIES_USAGE_REQ, today)
			}
		}
	}
	if incrReqKey != "" {
		_, err := gredis.Incrby(incrReqKey, length, 0)
		return err
	}
	return nil
}

func (u *UserToken) CollectRespUsage(id int, piID int, isStatic int, dataCenter int, length int) error {
	proxyIPService := proxy_ip_service.ProxyIP{ID: piID}
	proxyIP, err := proxyIPService.GetByID()
	if err != nil {
		return nil
	}
	var incrRespKey = ""
	today := coarsetime.CeilingTimezoneTimeNowYYMMDD(0)
	if isStatic == 0 {
		// 动态代理
		if proxyIP.ForwardPort == 22225 { // luminati 线路2
			incrRespKey = fmt.Sprintf("%s_%d", e.CACHE_LUMINATI_RESIDENTIAL_USAGE_RESP, today)
		} else if proxyIP.ForwardPort == 7777 { // oxylab 线路1
			incrRespKey = fmt.Sprintf("%s_%d", e.CACHE_OXYLAB_RESIDENTIAL_USAGE_RESP, today)
		} else {
			incrRespKey = fmt.Sprintf("%s_%d", e.CACHE_922_RESIDENTIAL_USAGE_RESP, today)
		}
	} else if isStatic == 1 {
		// 静态代理
		if dataCenter == 0 {
			if proxyIP.ForwardPort == 22225 { // luminati
				incrRespKey = fmt.Sprintf("%s_%d", e.CACHE_LUMINATI_ISP_USAGE_RESP, today)
			} else { // iproyal
				incrRespKey = fmt.Sprintf("%s_%d", e.CACHE_IPROYAL_ISP_USAGE_RESP, today)
			}
		} else if dataCenter == 1 { //机房代理
			if proxyIP.ForwardPort == 22225 {
				incrRespKey = fmt.Sprintf("%s_%d", e.CACHE_LUMINATI_DATACENTER_USAGE_RESP, today)
			} else { // instant_proxy
				incrRespKey = fmt.Sprintf("%s_%d", e.CACHE_INSTANTPROXIES_USAGE_RESP, today)
			}
		}
	}
	if incrRespKey != "" {
		_, err := gredis.Incrby(incrRespKey, length, 0)
		return err
	}

	return nil
}
