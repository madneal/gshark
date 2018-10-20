package models

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestGetCodeResultDetailById(t *testing.T) {
	id := int64(321)
	detail, _ := GetCodeResultDetailById(id)
	fmt.Println(detail)
}

func TestGetMatchedTests(t *testing.T) {
	_, codeResult, _ := GetReportById(321, true)
	var texts []*string
	texts = getMatchedTests(codeResult)
	assert.True(t, strings.Contains(*texts[0], "spdb"))
}

func TestChangeReportsStatusByRepo(t *testing.T) {
	var codeResult CodeResult
	has, err := Engine.Table("code_result").Cols("id", "status").Get(&codeResult)
	if err == nil && has {
		id := codeResult.Id
		var repo string
		has1, err1 := Engine.Table("code_result").Where("id=?", id).
			Cols("repo_name").Get(&repo)
		if err1 == nil && has1 {
			ChangeReportsStatusByRepo(id, 2)

			Engine.Table("code_result").Where("repo_name=?", repo).Get(&codeResult)
			assert.Equal(t, 2, codeResult.Status, "the status should be updated to 2")
		}
	}

}
