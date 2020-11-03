package models

type ProxyMachine struct {
	Model
	IpID int64 `json:"ip_id"`
}

func (ProxyMachine) TableName() string {
	return "t_proxy_machines"
}
