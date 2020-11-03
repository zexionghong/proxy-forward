package proxy_machine_service

import (
	"encoding/json"
	"proxy-forward/internal/models"
	"proxy-forward/internal/service/cache_service"
	"proxy-forward/pkg/gredis"
	"proxy-forward/pkg/logging"
)

type ProxyMachine struct {
	ID int
}

func (p *ProxyMachine) GetByID() (*models.ProxyMachine, error) {
	var cacheProxyMachine *models.ProxyMachine
	cache := cache_service.ProxyMachine{ID: p.ID}
	key := cache.GetKey()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Log.Warn(err.Error())
		} else {
			json.Unmarshal(data, &cacheProxyMachine)
			return cacheProxyMachine, nil
		}
	}

	proxyMachine, err := models.GetProxyMachineByID(p.ID)
	if err != nil {
		return nil, err
	}
	gredis.Set(key, proxyMachine, 60)
	return proxyMachine, nil
}
