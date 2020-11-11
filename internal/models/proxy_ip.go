package models

type ProxyIP struct {
	Model
	PmID        int `json:"pm_id"`
	ForwardPort int `json:"forward_port"`
	Port        int `json:"port"`
	IpID        int `json:"ip_id"`
	Online      int `json:"online"`
}

func (ProxyIP) TableName() string {
	return "t_proxy_ips"
}

// GetProxyIPByID get a single proxy_ip based on id
func GetProxyIPByID(id int) (*ProxyIP, error) {
	var result ProxyIP
	if err := db.Where("id = ? and deleted_on = ?", id, 0).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
