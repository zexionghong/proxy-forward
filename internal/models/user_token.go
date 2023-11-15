package models

import (
	"fmt"
	"strconv"
)

type UserToken struct {
	Model
	Username        string `json:"username" gorm:"column:username"`
	PassWord        string `json:"password" gorm:"column:password"`
	IsDeleted       bool   `json:"is_deleted" gorm:"column:is_deleted"`
	Expired         bool   `json:"is_expired" gorm:"column:is_expired"`
	IsStatic        bool   `json:"is_static" gorm:"column:is_static"`
	IsTraffic       bool   `json:"is_traffic" gorm:"column:is_traffic"`
	DataCenter      bool   `json:"is_datacenter" gorm:"column:is_datacenter"`
	PiID            int    `json:"proxy_id" gorm:"column:proxy_id"`
	Uid             int    `json:"user_id" gorm:"column:user_id"`
	Status          int    `json:"status" gorm:"column:status"`
	IsUse           bool   `json:"is_use" gorm:"column:is_use"`
	Host            string `json:"host" gorm:"column:host"`
	Port            int    `json:"port" gorm:"column:port"`
	ExpiredOn       int    `json:"expired_on" gorm:"column:expired_on"`
	ReqUsageAmount  int    `json:"request_amount" gorm:"column:request_amount"`
	RespUsageAmount int    `json:"response_amount" gorm:"column:response_amount"`
	Suffix          string
}

//
//type UserToken struct {
//	Model
//	Uid             int    `json:"uid"`
//	Username        string `json:"username"`
//	Passwd          string `json:"passwd"`
//	PiID            int    `json:"pi_id"`
//	Requests        int    `json:"requests"`
//	Traffic         int    `json:"traffic"`
//	Expired         int    `json:"expired"`
//	ReqUsageAmount  int    `json:"req_usage_amount"`
//	RespUsageAmount int    `json:"resp_usage_amount"`
//	PaywayID        int    `json:"payway_id"`
//	IsApi           int    `json:"is_api"`
//	LaID            int    `json:"la_id"`
//	PsID            int    `json:"ps_id"`
//	IsDeleted       int    `json:"is_deleted"`
//	IsStatic        int    `json:"is_static"`
//	DataCenter      int    `json:"data_center"`
//	Suffix          string
//}

func (UserToken) TableName() string {
	return "t_user_tokens"
}

// GetUserToken Get a single user_token based on username and passwd
func GetUserToken(username, passwd string) (*UserToken, error) {
	var result UserToken
	if err := db.Where("username = ? and password = ?  and deleted_on = ?", username, passwd, 0).First(&result).Error; err != nil {
		return nil, err
	}
	fmt.Print(strconv.Itoa(result.Uid))

	return &result, nil
}
