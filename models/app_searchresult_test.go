package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppSearchResult_Exist(t *testing.T) {
	var name = "借条贷"
	var market = "HUAWEI"
	appSearchReult := AppSearchResult{
		Name:   &name,
		Market: &market,
	}
	has, _ := appSearchReult.Exist()
	assert.True(t, has, "The result should exist!!!")
}

func TestChangeAppResultStatus(t *testing.T) {
	id := int64(1)
	err := changeAppResultStatus(1, id)
	assert.True(t, err == nil)
	appSearchResult := new(AppSearchResult)
	has, err := Engine.Id(id).Get(appSearchResult)
	if has {
		assert.True(t, appSearchResult.Status == 1)
	}
}
