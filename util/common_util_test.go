package util

import (
	"testing"
	"fmt"
	"../util"
)

func TestGetPageList(t *testing.T) {
	p := 30
	pages := 100
	//pageList := make([]int, 0)
	pageList := util.GetPageList(p, pages)
	fmt.Println(pageList)
}
