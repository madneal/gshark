package utils_test

import (
	"os"
	"testing"

	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/initialize"
	"github.com/madneal/gshark/utils"
)

const externalNotificationTestEnv = "GSHARK_NOTIFY_TESTS"

func requireExternalNotificationTest(t *testing.T) {
	t.Helper()
	if os.Getenv(externalNotificationTestEnv) != "1" {
		t.Skipf("set %s=1 to run external notification tests", externalNotificationTestEnv)
	}
}

func TestEmailSendExternal(t *testing.T) {
	requireExternalNotificationTest(t)
	global.GVA_VP = initialize.Viper("/Users/neal/project/gshark/server/config.yaml") // 初始化Viper
	if err := utils.EmailSend("test", "test"); err != nil {
		t.Fatal(err)
	}
}

func TestBotSendExternal(t *testing.T) {
	requireExternalNotificationTest(t)
	global.GVA_VP = initialize.Viper("/Users/neal/project/gshark/server/config.yaml")
	if err := utils.BotSend("Github敏感信息报告\n" + "test"); err != nil {
		t.Fatal(err)
	}
}
