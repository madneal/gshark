package routers

import (
	"testing"
	"fmt"
)

func TestGenerateFileHash(t *testing.T) {
	filepath := "app_assets.go"
	fmt.Println(GenerateFileHash(filepath))
}

