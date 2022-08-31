package bypass_service

import (
	"encoding/json"
	"proxy-forward/internal/models"
	"proxy-forward/internal/service/cache_service"
	"proxy-forward/pkg/gredis"
	"proxy-forward/pkg/logging"
)

type Bypass struct {
	PsID int
}

func (b *Bypass) GetBypassByPsID() ([]string, error) {
	var result []string
	cache := cache_service.Bypass{PsID: b.PsID}
	key := cache.GetKeyByPsID()
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Log.Warn(err.Error())
		} else {
			json.Unmarshal(data, &result)
			return result, nil
		}
	}
	bypass, err := models.GetBypassByPsID(b.PsID)
	if err != nil {
		return nil, err
	}
	if len(bypass) > 0 {
		for _, item := range bypass {
			result = append(result, item.Domain)
		}
	}
	gredis.Set(key, result, 600)
	return result, nil
}
