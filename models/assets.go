package models

import "time"

type Asset struct {
	Id          int64
	Owner       string `xorm:"text"`
	Type        string `xorm:"text notnull"`
	Content     string `xorm:"text notnull"`
	Status      int
	CreatedTime time.Time `xorm:"created"`
	UpdatedTime time.Time `xorm:"updated"`
}
