package proxy_supplier_service

import (
	"encoding/json"
	"proxy-forward/internal/models"
	"proxy-forward/internal/service/cache_service"
	"proxy-forward/pkg/gredis"
	"proxy-forward/pkg/logging"
)

type ProxySupplier struct {
	ID int
}

func (p *ProxySupplier) GetByID() (*models.ProxySupplier, error) {
	var cacheProxySupplier *models.ProxySupplier
	cache := cache_service.ProxySupplier{ID: p.ID}
	key := cache.GetKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Log.Warn(err.Error())
		} else {
			json.Unmarshal(data, &cacheProxySupplier)
			return cacheProxySupplier, nil
		}
	}

	proxySupplier, err := models.GetProxySupplierByID(p.ID)
	if err != nil {
		return nil, err
	}
	gredis.Set(key, proxySupplier, 60)
	return proxySupplier, nil
}
