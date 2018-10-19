package routers

import (
	"testing"
	"x-patrol/models"
	"fmt"
)

func TestSetUserInfoOfCodeResultDetail(t *testing.T)  {
	codeResultDetail, _ := models.GetCodeResultDetailById(int64(1))
	setUserInfoOfCodeResultDetail(codeResultDetail)
	fmt.Println(codeResultDetail.Blog)
}
