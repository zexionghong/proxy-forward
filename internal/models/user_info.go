package models

import (
	"encoding/json"
	"fmt"
	"proxy-forward/pkg/gredis"
	"proxy-forward/pkg/logging"
)

type UserInfo struct {
	Uid     int `json:"id"`
	Balance int `json:"balance"`
}

func (UserInfo) TableName() string {
	return "t_user_info"
}

func CanUse(userId int) (bool, error) {
	var result UserInfo
	key := fmt.Sprintf("USER_BALANCE_%d", userId)
	if gredis.Exists(key) {
		data, err := gredis.Get(key)
		if err != nil {
			logging.Log.Warn(err.Error())
		} else {
			json.Unmarshal(data, &result)
			return result.Balance >= 0, nil
		}
	}
	if err := db.Where("user_id = ?", userId).First(&result).Error; err != nil {
		return false, err
	}
	gredis.Set(key, result, 3600)
	if result.Balance >= 0 {
		return true, nil
	} else {
		return false, nil
	}
}
