package proxy_ip_service

import (
	"encoding/json"
	"proxy-forward/internal/models"
	"proxy-forward/internal/service/cache_service"
	"proxy-forward/pkg/gredis"
	"proxy-forward/pkg/logging"
)

type ProxyIP struct {
	ID int
}

func (p *ProxyIP) GetByID() (*models.ProxyIP, error) {
	var cacheProxyIP *models.ProxyIP
	cache := cache_service.ProxyIP{ID: p.ID}
	key := cache.GetKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Log.Warn(err.Error())
		} else {
			json.Unmarshal(data, &cacheProxyIP)
			return cacheProxyIP, nil
		}
	}
	proxyIP, err := models.GetProxyIPByID(p.ID)
	if err != nil {
		return nil, err
	}
	gredis.Set(key, proxyIP, 120)
	return proxyIP, nil
}

func (p *ProxyIP) DelteCache() (bool, error) {
	cache := cache_service.ProxyIP{ID: p.ID}
	key := cache.GetKey()
	if gredis.Exists(key) {
		return gredis.Delete(key)
	}
	return true, nil
}
