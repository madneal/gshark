// 自动生成模板Repo
package model

import (
	"github.com/madneal/gshark/global"
      "time"
)

// 如果含有time.Time 请自行import time包
type Repo struct {
      global.GVA_MODEL
      ProjectId  int `json:"projectId" form:"projectId" gorm:"column:project_id;comment:;type:int;"`
      Type  string `json:"type" form:"type" gorm:"column:type;comment:;type:varchar(10);size:10;"`
      Desc  string `json:"desc" form:"desc" gorm:"column:desc;comment:;type:varchar(300);size:300;"`
      Url  string `json:"url" form:"url" gorm:"column:url;comment:;type:varchar(300);size:300;"`
      Path  string `json:"path" form:"path" gorm:"column:path;comment:;type:varchar(100);size:100;"`
      LastActivityAt  time.Time `json:"lastActivityAt" form:"lastActivityAt" gorm:"column:last_activity_at;comment:"`
      Status  int `json:"status" form:"status" gorm:"column:status;comment:"`
}


func (Repo) TableName() string {
  return "repo"
}

