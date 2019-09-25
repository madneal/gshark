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

func TestCodeResult_Exist(t *testing.T) {
	url := "034397bb"
	codeResult := CodeResult{
		HTMLURL: &url,
	}
	exist, err := codeResult.Exist()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(exist)
}

//func TestChangeReportsStatusByRepo(t *testing.T) {
//	var codeResult CodeResult
//	has, err := Engine.Table("code_result").Cols("id", "status").Get(&codeResult)
//	if err == nil && has {
//		id := codeResult.Id
//		var repo string
//		has1, err1 := Engine.Table("code_result").Where("id=?", id).
//			Cols("repo_name").Get(&repo)
//		if err1 == nil && has1 {
//			ChangeReportsStatusByRepo(id, 2)
//
//			Engine.Table("code_result").Where("repo_name=?", repo).Get(&codeResult)
//			assert.Equal(t, )
//		}
//	}
//
//}

func TestGetMatchesTexts(t *testing.T) {
	repoName := "ralf-yin/lms"
	textMatches := GetMatchedTexts(repoName)
	for index, text := range textMatches {
		fmt.Println(index)
		fmt.Println(*text.Fragment)
	}
}
