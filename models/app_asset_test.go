package models

import (
	"testing"
	"fmt"
)

func TestDetect(t *testing.T) {
	var hash = "123456"
	var name = "app"
	appAsset := AppAsset {
		Name: &name,
		Hash: &hash,
	}
	num, err := appAsset.Insert()
	fmt.Println(num)
	if err != nil {
		fmt.Println(err)
	}
	isExist := Detect(hash)
	fmt.Println(isExist)
}
