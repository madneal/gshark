package models

import (
	"testing"
	"fmt"
)

func TestGetCodeResultDetailById(t *testing.T) {
	id := int64(321)
	detail, _ := GetCodeResultDetailById(id)
	fmt.Println(detail)
}
