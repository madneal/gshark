package routers

import (
	"gshark/models"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSetUserInfoOfCodeResultDetail(t *testing.T) {
	codeResultDetail, _ := models.GetCodeResultDetailById(int64(1))
	setUserInfoOfCodeResultDetail(codeResultDetail)
	assert.Equal(t, *codeResultDetail.Blog, "http://networkx.github.io/")
}
