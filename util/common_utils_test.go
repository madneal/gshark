package util

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

func TestSendMessage(t *testing.T) {
	SendMessage("SCU66439Tb628153bb4f665087db3ba7673d9b5cb5dcd6791379a1", "test", "test111")
}
