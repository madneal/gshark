package utils_test

import (
	"fmt"
	"github.com/madneal/gshark/core"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/utils"

	"testing"
)

func TestEmailTest(t *testing.T) {
	global.GVA_VP = core.Viper("/Users/neal/project/gshark/server/config.yaml") // 初始化Viper
	err := utils.EmailSend("test", "test")
	if err != nil {
		fmt.Print(err)
	}
}
