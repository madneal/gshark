package models

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAppSearchResult_Exist(t *testing.T) {
	var name = "借条贷"
	var market = "HUAWEI"
	appSearchReult := AppSearchResult{
		Name: &name,
		Market: &market,
	}
	has, _ := appSearchReult.Exist()
	assert.True(t, has, "The result should exist!!!")
}
