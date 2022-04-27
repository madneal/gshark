package gobuster

import (
	"context"
	"testing"
)

func TestRunDNS(t *testing.T) {
	RunDNS1(context.Background(), "baidu.com")
}
