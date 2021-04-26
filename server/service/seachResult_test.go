package service

import (
	"fmt"
	"github.com/madneal/gshark/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckExistOfSearchResult(t *testing.T) {
	result := model.SearchResult{
		Url: "https://baidu.com",
	}
	err, exist := CheckExistOfSearchResult(&result)
	assert.Equal(t, true, exist, "it should exits")
	if err != nil {
		fmt.Println(err)
	}
}
