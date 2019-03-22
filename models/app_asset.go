package models

import (
	"time"
	"github.com/neal1991/gshark/vars"
	"fmt"
)

type AppAsset struct {
	Id          int64
	Name        *string `json:"name,omitempty"`
	Description *string
	Market      *string `json:"market,omitempty"`
	Developer   *string
	Version     *string
	DeployDate  *string
	Url         *string
	// obtain from virustotal
	Sha256      *string
	Status       int
	CreatedTime time.Time `xorm:"created"`
	UpdatedTime time.Time `xorm:"updated"`
}

// sha256 is utilized to detect if the app exists
func Detect(sha256 string) (bool) {
	result, err := Engine.Table("app_assets").Where("sha256=?", sha256).Exist()
	if err != nil {
		fmt.Println(err)
	}
	return result
}


func ListAppAssets(page int) ([]InputInfo, int, error) {
	inputs := make([]InputInfo, 0)

	totalPages, err := Engine.Table("app_assets").Count()
	var pages int

	if int(totalPages)%vars.PAGE_SIZE == 0 {
		pages = int(totalPages) / vars.PAGE_SIZE
	} else {
		pages = int(totalPages)/vars.PAGE_SIZE + 1
	}

	if page >= pages {
		page = pages
	}

	if page < 1 {
		page = 1
	}

	err = Engine.Table("input_info").Limit(vars.PAGE_SIZE, (page-1)*vars.PAGE_SIZE).Find(&inputs)

	return inputs, pages, err
}

