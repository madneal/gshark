package routers

import (
	"testing"
	"gshark/models"
	"fmt"
)

func TestSetUserInfoOfCodeResultDetail(t *testing.T)  {
	codeResultDetail, _ := models.GetCodeResultDetailById(int64(1))
	setUserInfoOfCodeResultDetail(codeResultDetail)
	fmt.Println(codeResultDetail.Blog)
}
