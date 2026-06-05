package model

type CasbinModel struct {
	Ptype       string `json:"ptype" gorm:"column:p_type;size:100"`
	AuthorityId string `json:"rolename" gorm:"column:v0;size:100"`
	Path        string `json:"path" gorm:"column:v1;size:100"`
	Method      string `json:"method" gorm:"column:v2;size:100"`
	V3          string `json:"v3" gorm:"column:v3;size:100"`
	V4          string `json:"v4" gorm:"column:v4;size:100"`
	V5          string `json:"v5" gorm:"column:v5;size:100"`
}

func (CasbinModel) TableName() string {
	return "casbin_rule"
}
