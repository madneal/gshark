package service

import (
	"github.com/madneal/gshark/config"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/utils"
	"go.uber.org/zap"
)

func GetSystemConfig() (err error, conf config.Server) {
	return nil, global.GVA_CONFIG
}

func SetSystemConfig(system model.System) (err error) {
	cs := utils.StructToMap(system.Config)
	for k, v := range cs {
		global.GVA_VP.Set(k, v)
	}
	if err = global.GVA_VP.WriteConfig(); err != nil {
		return err
	}
	// WriteConfig persists the values but does not refresh the runtime snapshot
	// consumed by the rest of the application. Keep both sources in sync so a
	// successful update is immediately visible to GetSystemConfig and services.
	return global.GVA_VP.Unmarshal(&global.GVA_CONFIG)
}

func GetServerInfo() (server *utils.Server, err error) {
	var s utils.Server
	s.Os = utils.InitOS()
	if s.Cpu, err = utils.InitCPU(); err != nil {
		global.GVA_LOG.Error("func utils.InitCPU() Failed!", zap.String("err", err.Error()))
		return &s, err
	}
	if s.Rrm, err = utils.InitRAM(); err != nil {
		global.GVA_LOG.Error("func utils.InitRAM() Failed!", zap.String("err", err.Error()))
		return &s, err
	}
	if s.Disk, err = utils.InitDisk(); err != nil {
		global.GVA_LOG.Error("func utils.InitDisk() Failed!", zap.String("err", err.Error()))
		return &s, err
	}

	return &s, nil
}
