package models

type UserToken struct {
	Model
	Uid      int    `json:"uid"`
	Username string `json:"username"`
	Passwd   string `json:"passwd"`
	PiID     int    `json:"pi_id"`
	Requests int    `json:"requests"`
	Traffic  int    `json:"traffic"`
}

func (UserToken) TableName() string {
	return "t_user_tokens"
}

// GetUserToken Get a single user_token based on username and passwd
func GetUserToken(username, passwd string) (*UserToken, error) {
	var result UserToken
	if err := db.Where("username = ? and passwd = ? and deleted_on = ?", username, passwd, 0).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}