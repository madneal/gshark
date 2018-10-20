package routers

import (
	"fmt"
	"gshark/models"
	"testing"
)

func TestSetUserInfoOfCodeResultDetail(t *testing.T) {
	codeResultDetail, _ := models.GetCodeResultDetailById(int64(1))
	setUserInfoOfCodeResultDetail(codeResultDetail)
	fmt.Println(codeResultDetail.Blog)
}
