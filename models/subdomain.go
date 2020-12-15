package models

import (
	"github.com/madneal/gshark/logger"
	"github.com/madneal/gshark/util"
	"github.com/madneal/gshark/vars"
	"time"
)

type Subdomain struct {
	Id        int64
	Domain    *string
	Status    int
	Subdomain *string `xorm:"unique"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (subdomain *Subdomain) Insert() (int64, error) {
	return Engine.Insert(subdomain)
}

func ListSubdomains() ([]Subdomain, error) {
	subdomains := make([]Subdomain, 0)
	err := Engine.Table("subdomain").Find(&subdomains)
	return subdomains, err
}

func ListSubdomainsByPage(page int) ([]Subdomain, int, int) {
	results := make([]Subdomain, 0)
	count, err := Engine.Table("subdomain").Where("status = 0").Count()
	err = Engine.Table("subdomain").Where("status = 0").Limit(vars.PAGE_SIZE, (page-1)*vars.PAGE_SIZE).
		Desc("id").Find(&results)
	if err != nil {
		logger.Log.Error(err)
	}
	page, pages := util.GetPageAndPagesByCount(page, int(count))
	return results, pages, int(count)
}

func IgnoreSubdomain(id int) error {
	subdomain := new(Subdomain)
	subdomain.Status = 1
	_, err := Engine.ID(int64(id)).Update(subdomain)
	return err
}
