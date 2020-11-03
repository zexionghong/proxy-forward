package models

type IP struct {
	Model
	IpAddr int64 `json:"ip"`
}

func (IP) TableName() string {
	return "t_ips"
}

// GetIPByID get a single ip based on id
func GetIPByID(id int) (*IP, error) {
	var result IP
	if err := db.Where("id = ? and deleted_on = ?", id, 0).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
