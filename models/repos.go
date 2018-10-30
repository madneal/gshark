package models

import (
	"gshark/vars"
)

type Repo struct {
	Id     int64
	Name   string
	Url    string
	Src    *InputInfo
	Status int `xorm:"int notnull default(1)"`
}

func NewRepo(name, repoUrl string, src *InputInfo) (repo *Repo) {
	return &Repo{Name: name, Url: repoUrl, Src: src, Status: 1}
}

func (r *Repo) Insert() (int64, error) {
	return Engine.Insert(r)
}

func (r *Repo) Exist() (bool, error) {
	repo := new(Repo)
	repo.Name = r.Name
	return Engine.Table("repo").Get(repo)
}

func ListReposPage(page int) ([]Repo, int, error) {
	repos := make([]Repo, 0)
	totalPages, err := Engine.Table("repo").Count()
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

	err = Engine.Limit(vars.PAGE_SIZE, (page-1)*vars.PAGE_SIZE).Find(&repos)
	return repos, pages, err
}

func ListEnableRepos() ([]Repo, error) {
	repos := make([]Repo, 0)
	err := Engine.Table("repo").Where("status=?", 1).Find(&repos)
	return repos, err
}

func EnableRepoById(id int64) error {
	repo := new(Repo)
	has, err := Engine.ID(id).Get(repo)
	if err == nil && has {
		repo.Status = 1
		_, err = Engine.ID(id).Cols("status").Update(repo)
	}
	return err
}

func DisableRepoById(id int64) error {
	repo := new(Repo)
	has, err := Engine.ID(id).Get(repo)
	if err == nil && has {
		repo.Status = 0
		_, err = Engine.ID(id).Cols("status").Update(repo)
	}
	return err
}

func DeleteAllRepos() error {
	sqlCmd := "delete from repo"
	_, err := Engine.Exec(sqlCmd)
	return err
}

func DisableRepoByUrl(repoUrl string) error {
	repo := new(Repo)
	has, err := Engine.Table("repo").Where("url=?", repoUrl).Get(repo)
	if err == nil && has {
		repo.Status = 0
		_, err = Engine.ID(repo.Id).Cols("status").Update(repo)
	}
	return err
}
