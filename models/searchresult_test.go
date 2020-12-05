package models

import (
	"fmt"
	"github.com/madneal/gshark/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetCodeResultDetailById(t *testing.T) {
	id := int64(321)
	detail, _ := GetCodeResultDetailById(id)
	fmt.Println(detail)
}

func TestCodeResult_Exist(t *testing.T) {
	codeResult := new(CodeResult)
	f := "                                    \\u003cspan class=\\\"views\\\"\\u003e???\\u003cb style=\\\"color: red;\\\"\\u003e64220\\u003c/b\\u003e\\u003c/span\\u003e\\n                                    \\u003cspan class=\\\"name\\\"\\u003e\\u003ca href=\\\"http://www.meituan.com/r/i1186336\\\" target=\\\"_blank\\\"\\u003e??\\u003c/a\\u003e\\u003c/span\\u003e\\n"
	hash := util.GenMd5WithSpecificLen(f, 50)
	codeResult.Textmatchmd5 = &hash
	html := "www.baidu.com"
	codeResult.HTMLURL = &html
	exist, _ := codeResult.Exist()
	assert.Equal(t, true, exist, "the result should exist")
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
