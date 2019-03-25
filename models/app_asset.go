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
	Filename    *string
	// obtain from virustotal
	Hash      *string
	Status       int
	CreatedTime time.Time `xorm:"created"`
	UpdatedTime time.Time `xorm:"updated"`
}

func NewAppAsset(name, desc, market, developer, version, deployDate, url, hash, filename string, status int) *AppAsset {
	appAsset := AppAsset{
		Name: &name,
		Description: &desc,
		Market: &market,
		Developer: &developer,
		Version: &version,
		DeployDate: &deployDate,
		Url: &url,
		Hash: &hash,
		Filename: &filename,
		Status: status,
	}
	appAsset.CreatedTime = time.Now().Local()
	return &appAsset
}

func Detect(hash string) (bool, int64) {
	appAsset := new(AppAsset)
	var id int64
	has, err := Engine.Table("app_asset").Where("hash=?", hash).Get(appAsset)
	if err != nil {
		fmt.Println(err)
	}
	if !has {
		id = -1
	} else {
		id = appAsset.Id
	}
	return has, id
}

func GetAppAssetById(id int64)  AppAsset {
	appAsset := new(AppAsset)
	_, err := Engine.Table("app_asset").Where("id=?", id).Get(appAsset)
	if err != nil {
		fmt.Println(err)
	}
	return *appAsset
}

func EditAppAssetById(id int64, asset *AppAsset) error {
	appAsset := new(AppAsset)
	has, err := Engine.Id(id).Get(appAsset)
	if err == nil && has {
		Engine.Id(id).Update(asset)
	}
	return err
}

func (r *AppAsset) Insert() (int64, error) {
	return Engine.Insert(r)
}

func DeleteAppAssetById(id int64) {
	appAsset := new(AppAsset)
	_, err := Engine.Id(id).Delete(appAsset)
	if err != nil {
		fmt.Println(err)
	}
}


func ListAppAssets(page int) ([]AppAsset, int, error) {
	apps := make([]AppAsset, 0)

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

	err = Engine.Table("app_asset").Limit(vars.PAGE_SIZE, (page-1)*vars.PAGE_SIZE).Find(&apps)

	return apps, pages, err
}

