package models

type User struct {
	Uid        int  `json:"id"`
	IsUse      bool `json:"is_use"`
	CreatedOn  int  `json:"created_on"`
	ModifiedOn int  `json:"updated_on" gorm:"column:updated_on"`
	DeletedOn  int  `json:"deleted_on"`
}

func (User) TableName() string {
	return "t_users"
}

// GetUser Get a single user based on uid
func GetUser(uid int) (*User, error) {
	var result User
	if err := db.Where("id = ?", uid).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
