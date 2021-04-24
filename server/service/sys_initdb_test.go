package service

import (
	"fmt"
	"github.com/madneal/gshark/model/request"
	"testing"
)

func TestRemoveDB(t *testing.T) {
	conf := request.InitDB{
		Host:     "localhost",
		DBName:   "gshark",
		UserName: "gshark",
		Password: "gshark",
	}
	err := RemoveDB(conf)
	if err != nil {
		fmt.Println(err)
	}
}
