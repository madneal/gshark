package service

import (
	"errors"
	"strings"

	"github.com/casbin/casbin/util"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
)

func UpdateCasbin(authorityId string, casbinInfos []request.CasbinInfo) error {
	ClearCasbin(0, authorityId)
	rules := [][]string{}
	for _, v := range casbinInfos {
		cm := model.CasbinModel{
			Ptype:       "p",
			AuthorityId: authorityId,
			Path:        v.Path,
			Method:      v.Method,
		}
		rules = append(rules, []string{cm.AuthorityId, cm.Path, cm.Method})
	}
	e, err := Casbin()
	if err != nil {
		return err
	}
	success, _ := e.AddPolicies(rules)
	if success == false {
		return errors.New("存在相同api,添加失败,请联系管理员")
	}
	return nil
}

func UpdateCasbinApi(oldPath string, newPath string, oldMethod string, newMethod string) error {
	err := global.GVA_DB.Table("casbin_rule").Model(&model.CasbinModel{}).Where("v1 = ? AND v2 = ?", oldPath, oldMethod).Updates(map[string]interface{}{
		"v1": newPath,
		"v2": newMethod,
	}).Error
	return err
}

func GetPolicyPathByAuthorityId(authorityId string) (pathMaps []request.CasbinInfo, err error) {
	e, err := Casbin()
	if err != nil {
		return nil, err
	}
	list := e.GetFilteredPolicy(0, authorityId)
	for _, v := range list {
		pathMaps = append(pathMaps, request.CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps, nil
}

func ClearCasbin(v int, p ...string) bool {
	e, err := Casbin()
	if err != nil {
		return false
	}
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success

}

func Casbin() (*casbin.Enforcer, error) {
	a, err := gormadapter.NewAdapterByDB(global.GVA_DB)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(global.GVA_CONFIG.Casbin.ModelPath, a)
	if err != nil {
		return nil, err
	}
	e.AddFunction("ParamsMatch", ParamsMatchFunc)
	if err := e.LoadPolicy(); err != nil {
		return nil, err
	}
	return e, nil
}

func ParamsMatch(fullNameKey1 string, key2 string) bool {
	key1 := strings.Split(fullNameKey1, "?")[0]
	// 剥离路径后再使用casbin的keyMatch2
	return util.KeyMatch2(key1, key2)
}

func ParamsMatchFunc(args ...interface{}) (interface{}, error) {
	name1 := args[0].(string)
	name2 := args[1].(string)

	return ParamsMatch(name1, name2), nil
}
