package util

import (
	"testing"
	"fmt"
	"../util"
)

func TestGetPageList(t *testing.T) {
	testPList := []int{1, 5, 9, 25, 100}
	pages := 100
	for _, p := range testPList {
		pageList := util.GetPageList(p, 5, pages)
		fmt.Printf("p:%d-- the length of pageList:%d  %v\n", p, len(pageList), pageList)
	}
}
