package models

import (
	"testing"
	"fmt"
	"github.com/neal1991/x-patrol/models"
)

func TestGetCodeResultDetailById(t *testing.T) {
	id := int64(321)
	detail, _ := models.GetCodeResultDetailById(id)
	fmt.Println(detail)
}
