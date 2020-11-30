package gobuster

import (
	"fmt"
	"testing"
	"time"
)

func TestRunDns(t *testing.T) {
	err := RunDNS("meituan.com")
	if err != nil {
		fmt.Println(err)
	}
}

func TestRunTask(t *testing.T) {
	RunTask(time.Duration(900))
}
