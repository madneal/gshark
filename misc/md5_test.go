package misc

import (
	"fmt"
	"testing"
)

func TestMakeMd5(t *testing.T) {
	pass := "gshark"
	fmt.Println(MakeMd5(pass))
}
