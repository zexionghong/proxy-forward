package models

type ProxyMachine struct {
	Model
	IpID int `json:"ip_id"`
	PsID int `json:"ps_id"`
}

func (ProxyMachine) TableName() string {
	return "t_proxy_machines"
}

// GetProxyMachineByID get a single proxy_machine based on id
func GetProxyMachineByID(id int) (*ProxyMachine, error) {
	var result ProxyMachine
	if err := db.Where("id = ? and deleted_on = ?", id, 0).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}
