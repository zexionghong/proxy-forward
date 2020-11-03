package ip_service

import (
	"encoding/json"
	"proxy-forward/internal/models"
	"proxy-forward/internal/service/cache_service"
	"proxy-forward/pkg/gredis"
	"proxy-forward/pkg/logging"
)

type IP struct {
	ID int
}

func (ip *IP) GetByID() (*models.IP, error) {
	var cacheIP *models.IP
	cache := cache_service.IP{ID: ip.ID}
	key := cache.GetKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Log.Warn(err.Error())
		} else {
			json.Unmarshal(data, &cacheIP)
			return cacheIP, nil
		}
	}

	iP, err := models.GetIPByID(ip.ID)
	if err != nil {
		return nil, err
	}
	gredis.Set(key, iP, 3600)
	return iP, nil
}
