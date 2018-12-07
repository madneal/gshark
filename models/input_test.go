package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInputInfo_Exist(t *testing.T) {
	repoUrl := "https://github.com/Yvoox/SPFL"
	fullName := "Yvoox/SPFL"
	inputInfo := NewInputInfo("repo", repoUrl, fullName)
	result, _ := inputInfo.Exist(repoUrl)
	assert.True(t, result)
}
