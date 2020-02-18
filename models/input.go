package models

import (
	"github.com/neal1991/gshark/vars"

	"time"
)

type InputInfo struct {
	Id          int64
	Type        string `xorm:"varchar(255) notnull"`
	Url         string `xorm:"text notnull"`
	Path        string `xorm:"text notnull"`
	Developer   string `xorm:"text"`
	Status      int
	ProjectId   int
	Version     int       `xorm:"version"`
	CreatedTime time.Time `xorm:"created"`
	UpdatedTime time.Time `xorm:"updated"`
}

func NewInputInfo(inputType, url, path string) (info *InputInfo) {
	return &InputInfo{Type: inputType, Url: url, Path: path}
}

func (i *InputInfo) Insert() (int64, error) {
	return Engine.Insert(i)
}

func (i *InputInfo) Exist() (bool, error) {
	return Engine.Exist(&InputInfo{
		Url: i.Url,
	})
}

func GetInputInfoById(id int64) (*InputInfo, bool, error) {
	input := InputInfo{Id: id}
	has, err := Engine.Table("input_info").ID(id).Get(&input)
	return &input, has, err
}

func EditInputInfoById(id int64, inputType, content, desc string) error {
	input := new(InputInfo)
	var err error
	has, err := Engine.ID(id).Get(input)
	if err == nil && has {
		input.Type = inputType
		input.Url = content
		input.Path = desc
		_, err = Engine.ID(id).Update(input)
	}
	return err
}

func UpdateStatusById(status, id int) error {
	_, err := Engine.Exec("update input_info set status = ? where project_id = ?", status, id)
	return err
}

func DeleteInputInfoById(id int64) error {
	input := new(InputInfo)
	_, err := Engine.Table("input_info").ID(id).Delete(input)
	return err
}

func (i InputInfo) DeleteByProjectId() error {
	_, err := Engine.Table("input_info").Where("project_id = ?", i.ProjectId).Delete(&i)
	return err
}

func DeleteAllInputInfo() error {
	sqlCMD := "delete from input_info;"
	_, err := Engine.Exec(sqlCMD)
	return err
}

func ListInputInfo() ([]InputInfo, error) {
	inputs := make([]InputInfo, 0)
	err := Engine.Table("input_info").Find(&inputs)
	return inputs, err
}

func ListInputInfoByType(infoType string) ([]InputInfo, error) {
	inputs := make([]InputInfo, 0)
	err := Engine.Table("input_info").Where("type = ?", infoType).Find(&inputs)
	return inputs, err
}

func ListInputInfoPage(page int) ([]InputInfo, int, error) {
	inputs := make([]InputInfo, 0)

	totalPages, err := Engine.Table("input_info").Count()
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
