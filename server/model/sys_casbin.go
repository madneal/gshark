package model

type CasbinModel struct {
	Ptype       string `json:"ptype" gorm:"column:p_type"`
	AuthorityId string `json:"rolename" gorm:"column:v0"`
	Path        string `json:"path" gorm:"column:v1"`
	Method      string `json:"method" gorm:"column:v2"`
	V3          string `json:"v3" gorm:"column:v3"`
	V4          string `json:"v4" gorm:"column:v4"`
	V5          string `json:"v5" gorm:"column:v5"`
}

func (CasbinModel) TableName() string {
	return "casbin_rule"
}
