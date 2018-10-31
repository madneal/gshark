package common

import (
	"fmt"
	"testing"
)

func TestGetPageList(t *testing.T) {
	testPList := []int{1, 5, 9, 25, 100}
	pages := 100
	for _, p := range testPList {
		pageList := GetPageList(p, 5, pages)
		fmt.Printf("p:%d-- the length of pageList:%d  %v\n", p, len(pageList), pageList)
	}
}
