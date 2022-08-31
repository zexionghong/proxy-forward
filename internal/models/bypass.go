package models

type Bypass struct {
	Model
	Domain string `json:"domain"`
}

func (Bypass) TableName() string {
	return "t_bypass"
}

func GetBypassByPsID(PsID int) ([]*Bypass, error) {
	var result []*Bypass
	err := db.Joins("LEFT JOIN t_bypass_proxy_suppliers ON t_bypass.id = t_bypass_proxy_suppliers.bypass_id ").Where("t_bypass_proxy_suppliers.ps_id = ?", PsID).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
