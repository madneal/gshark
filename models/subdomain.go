package models

import "time"

type Subdomain struct {
	Id        int64
	Domain    *string
	Subdomain *string `gorm:"index:subdomain_index,unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (subdomain *Subdomain) Insert() (int64, error) {
	return Engine.Insert(subdomain)
}
