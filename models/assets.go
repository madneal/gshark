package models

import "time"

type Asset struct {
	Id          int64
	Org         string `xorm:"text"`
	AssetType   string `xorm:"text notnull"`
	Content     string `xorm:"text notnull"`
	Status      int
	CreatedTime time.Time `xorm:"created"`
	UpdatedTime time.Time `xorm:"updated"`
}

func NewAsset(org, assetType, content string) *Asset {
	return &Asset{
		Org:         org,
		AssetType:   assetType,
		Content:     content,
		Status:      0,
		CreatedTime: time.Time{},
		UpdatedTime: time.Time{},
	}
}

func (asset *Asset) Insert() (int64, error) {
	return Engine.Insert(asset)
}

func (asset *Asset) Exists() (bool, error) {
	return Engine.Exist(&Asset{
		AssetType: asset.AssetType,
		Content:   asset.Content,
	})
}

func BatchInsert(assets []*Asset) error {
	for _, asset := range assets {
		exist, err := asset.Exists()
		if err != nil {
			return err
		}
		if exist == false {
			asset.Insert()
		}
	}
	return nil
}
