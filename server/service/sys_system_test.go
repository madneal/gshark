package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/madneal/gshark/config"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/spf13/viper"
)

func TestSetSystemConfigRefreshesRuntimeConfig(t *testing.T) {
	previousViper := global.GVA_VP
	previousConfig := global.GVA_CONFIG
	t.Cleanup(func() {
		global.GVA_VP = previousViper
		global.GVA_CONFIG = previousConfig
	})

	configPath := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(configPath, []byte("system:\n  addr: 8888\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	vp := viper.New()
	vp.SetConfigFile(configPath)
	if err := vp.ReadInConfig(); err != nil {
		t.Fatal(err)
	}
	global.GVA_VP = vp
	global.GVA_CONFIG = config.Server{System: config.System{Addr: 8888}}

	updated := global.GVA_CONFIG
	updated.System.Addr = 9999
	if err := SetSystemConfig(model.System{Config: updated}); err != nil {
		t.Fatal(err)
	}

	if got := global.GVA_CONFIG.System.Addr; got != 9999 {
		t.Fatalf("runtime config address = %d, want 9999", got)
	}
	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("config file was not persisted: %v", err)
	}
	if _, got := GetSystemConfig(); got.System.Addr != 9999 {
		t.Fatalf("returned config address = %d, want 9999", got.System.Addr)
	}
}
