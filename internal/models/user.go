package models

type User struct {
	Uid        int `json:"uid"`
	IsUse      int `json:"is_use"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

func (User) TableName() string {
	return "t_users"
}

// GetUser Get a single user based on uid
func GetUser(uid int) (*User, error) {
	var result User
	if err := db.Where("uid = ?", uid).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
