package misc

import (
	"testing"
	"fmt"
)

func TestMakeMd5(t *testing.T) {
	pass := "admin"
	fmt.Println(MakeMd5(pass))
}
