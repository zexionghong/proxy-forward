package models

type BypassProxySupplier struct {
	Model
	BypassID int `json:"bypass_id"`
	PsID     int `json:"ps_id"`
}

func (BypassProxySupplier) TableName() string {
	return "t_bypass_proxy_suppliers"
}
