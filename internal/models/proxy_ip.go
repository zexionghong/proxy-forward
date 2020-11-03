package models

type ProxyIP struct {
	Model
	PmID   int `json:"pm_id"`
	Port   int `json:"port"`
	IpID   int `json:"ip_id"`
	Online int `json:"online"`
}

func (ProxyIP) TableName() string {
	return "t_proxy_ips"
}
