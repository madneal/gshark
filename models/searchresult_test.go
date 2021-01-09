package models

import (
	"fmt"
	"testing"
)

func TestGetCodeResultDetailById(t *testing.T) {
	id := int64(321)
	detail, _ := GetCodeResultDetailById(id)
	fmt.Println(detail)
}

func TestGetMatchesTexts(t *testing.T) {
	repoName := "ralf-yin/lms"
	textMatches := GetMatchedTexts(repoName)
	for index, text := range textMatches {
		fmt.Println(index)
		fmt.Println(*text.Fragment)
	}
}
