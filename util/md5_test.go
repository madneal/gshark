package util

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeMd5(t *testing.T) {
	pass := "gshark"
	fmt.Println(MakeMd5(pass))
}

func TestGenMd5WithSpecificLen(t *testing.T) {
	f := "                                    \\u003cspan class=\\\"views\\\"\\u003e???\\u003cb style=\\\"color: red;\\\"\\u003e64220\\u003c/b\\u003e\\u003c/span\\u003e\\n                                    \\u003cspan class=\\\"name\\\"\\u003e\\u003ca href=\\\"http://www.meituan.com/r/i1186336\\\" target=\\\"_blank\\\"\\u003e??\\u003c/a\\u003e\\u003c/span\\u003e\\n"
	fmt.Println(GenMd5WithSpecificLen(f, 50))
	assert.Equal(t, true, GenMd5WithSpecificLen(f, 50) == "aa93d91535aae2a951d98c25c81c160a000000000000000000",
		"the genetated hash should be the same")
}
