package models

import (
	"github.com/madneal/gshark/logger"
	"testing"
)

func TestIgnoreSubdomain(t *testing.T) {
	err := IgnoreSubdomain(1)
	if err != nil {
		logger.Log.Error(err)
	}
}
