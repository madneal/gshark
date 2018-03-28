package util

import (
	"testing"
	"fmt"
	"../util"
)

func TestGetPageList(t *testing.T) {
	//p := 5
	testPList := []int{1, 5, 10, 23, 45, 100}
	pages := 100
	for _, p := range testPList {
		pageList := util.GetPageList(p, pages)
		fmt.Printf("p:%d-- the length of pageList:%d  %v\n", p, len(pageList), pageList)
	}
	//fmt.Println(pageList)
	//pageList := make([]int, 0)
}
