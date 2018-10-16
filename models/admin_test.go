package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	//fmt.Println(134)
	NewDbEngine()
	_, role, _ := Auth("admin", "123@admin")
	assert.Equal(t, "admin", role)
}
