package models

type ProxySupplier struct {
	Model
	Name      string `json:"name"`
	OnlyHttp  int    `json:"only_http"`
	SelfCheck int    `json:"self_check"`
}

func (ProxySupplier) TableName() string {
	return "t_proxy_suppliers"
}

// GetProxySupplierByID get a single proxy_supplier based on id
func GetProxySupplierByID(id int) (*ProxySupplier, error) {
	var result ProxySupplier
	if err := db.Where("id = ? and deleted_on = ?", id, 0).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
