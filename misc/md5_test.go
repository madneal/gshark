package misc

import (
	"testing"
	"fmt"
)

func TestMakeMd5(t *testing.T) {
	pass := "123456"
	fmt.Println(MakeMd5(pass))
}
