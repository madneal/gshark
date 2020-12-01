package gobuster

import (
	"testing"
	"time"
)

func TestRunDns(t *testing.T) {
	RunDNS("meituan.com")
}

func TestRunTask(t *testing.T) {
	RunTask(time.Duration(900))
}
