package gobuster

import (
	"fmt"
	"testing"
)

func TestRunDns(t *testing.T) {
	err := RunDNS("meituan.com")
	if err != nil {
		fmt.Println(err)
	}
}
