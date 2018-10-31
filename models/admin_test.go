package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	//fmt.Println(134)
	NewDbEngine()
	_, role, _ := Auth("admin", "admin")
	assert.Equal(t, "admin", role)
}
