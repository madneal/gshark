package routers

import (
	"github.com/madneal/gshark/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetUserInfoOfCodeResultDetail(t *testing.T) {
	codeResultDetail, _ := models.GetCodeResultDetailById(int64(1))
	setUserInfoOfCodeResultDetail(codeResultDetail)
	assert.Equal(t, *codeResultDetail.Blog, "http://networkx.github.io/")
}
