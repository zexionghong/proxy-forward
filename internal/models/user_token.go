package models

type UserToken struct {
	Model
	Uid             int    `json:"uid"`
	Username        string `json:"username"`
	Passwd          string `json:"passwd"`
	PiID            int    `json:"pi_id"`
	Requests        int    `json:"requests"`
	Traffic         int    `json:"traffic"`
	Expired         int    `json:"expired"`
	ReqUsageAmount  int    `json:"req_usage_amount"`
	RespUsageAmount int    `json:"resp_usage_amount"`
	PaywayID        int    `json:"payway_id"`
	IsApi           int    `json:"is_api"`
	Suffix          string
}

func (UserToken) TableName() string {
	return "t_user_tokens"
}

// GetUserToken Get a single user_token based on username and passwd
func GetUserToken(username, passwd string) (*UserToken, error) {
	var result UserToken
	if err := db.Where("username = ? and passwd = ? and payway_id = ? and deleted_on = ?", username, passwd, 0, 0).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
