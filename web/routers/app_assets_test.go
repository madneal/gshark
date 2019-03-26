package routers

import (
	"fmt"
	"testing"
)

func TestGenerateFileHash(t *testing.T) {
	filepath := "app_assets.go"
	fmt.Println(GenerateFileHash(filepath))
}
