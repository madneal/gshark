package source

import (
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gookit/color"
	"github.com/madneal/gshark/global"
	"gorm.io/gorm"
)

var Casbin = new(casbin)

type casbin struct{}

var carbines = []gormadapter.CasbinRule{
	{PType: "p", V0: "888", V1: "/base/login", V2: "POST"},
	{PType: "p", V0: "888", V1: "/user/register", V2: "POST"},
	{PType: "p", V0: "888", V1: "/api/createApi", V2: "POST"},
	{PType: "p", V0: "888", V1: "/api/getApiList", V2: "POST"},
	{PType: "p", V0: "888", V1: "/api/getApiById", V2: "POST"},
	{PType: "p", V0: "888", V1: "/api/deleteApi", V2: "POST"},
	{PType: "p", V0: "888", V1: "/api/updateApi", V2: "POST"},
	{PType: "p", V0: "888", V1: "/api/getAllApis", V2: "POST"},
	{PType: "p", V0: "888", V1: "/authority/createAuthority", V2: "POST"},
	{PType: "p", V0: "888", V1: "/authority/deleteAuthority", V2: "POST"},
	{PType: "p", V0: "888", V1: "/authority/getAuthorityList", V2: "POST"},
	{PType: "p", V0: "888", V1: "/authority/setDataAuthority", V2: "POST"},
	{PType: "p", V0: "888", V1: "/authority/updateAuthority", V2: "PUT"},
	{PType: "p", V0: "888", V1: "/authority/copyAuthority", V2: "POST"},
	{PType: "p", V0: "888", V1: "/menu/getMenu", V2: "POST"},
	{PType: "p", V0: "888", V1: "/menu/getMenuList", V2: "POST"},
	{PType: "p", V0: "888", V1: "/menu/addBaseMenu", V2: "POST"},
	{PType: "p", V0: "888", V1: "/menu/getBaseMenuTree", V2: "POST"},
	{PType: "p", V0: "888", V1: "/menu/addMenuAuthority", V2: "POST"},
	{PType: "p", V0: "888", V1: "/menu/getMenuAuthority", V2: "POST"},
	{PType: "p", V0: "888", V1: "/menu/deleteBaseMenu", V2: "POST"},
	{PType: "p", V0: "888", V1: "/menu/updateBaseMenu", V2: "POST"},
	{PType: "p", V0: "888", V1: "/menu/getBaseMenuById", V2: "POST"},
	{PType: "p", V0: "888", V1: "/user/changePassword", V2: "POST"},
	{PType: "p", V0: "888", V1: "/user/getUserList", V2: "POST"},
	{PType: "p", V0: "888", V1: "/user/setUserAuthority", V2: "POST"},
	{PType: "p", V0: "888", V1: "/user/deleteUser", V2: "DELETE"},
	{PType: "p", V0: "888", V1: "/casbin/updateCasbin", V2: "POST"},
	{PType: "p", V0: "888", V1: "/casbin/getPolicyPathByAuthorityId", V2: "POST"},
	{PType: "p", V0: "888", V1: "/casbin/casbinTest/:pathParam", V2: "GET"},
	{PType: "p", V0: "888", V1: "/jwt/jsonInBlacklist", V2: "POST"},
	{PType: "p", V0: "888", V1: "/system/getSystemConfig", V2: "POST"},
	{PType: "p", V0: "888", V1: "/system/setSystemConfig", V2: "POST"},
	{PType: "p", V0: "888", V1: "/system/getServerInfo", V2: "POST"},
	{PType: "p", V0: "888", V1: "/system/emailTest", V2: "POST"},
	{PType: "p", V0: "888", V1: "/system/botTest", V2: "GET"},
	{PType: "p", V0: "888", V1: "/sysDictionaryDetail/createSysDictionaryDetail", V2: "POST"},
	{PType: "p", V0: "888", V1: "/sysDictionaryDetail/deleteSysDictionaryDetail", V2: "DELETE"},
	{PType: "p", V0: "888", V1: "/sysDictionaryDetail/updateSysDictionaryDetail", V2: "PUT"},
	{PType: "p", V0: "888", V1: "/sysDictionaryDetail/findSysDictionaryDetail", V2: "GET"},
	{PType: "p", V0: "888", V1: "/sysDictionaryDetail/getSysDictionaryDetailList", V2: "GET"},
	{PType: "p", V0: "888", V1: "/sysDictionary/createSysDictionary", V2: "POST"},
	{PType: "p", V0: "888", V1: "/sysDictionary/deleteSysDictionary", V2: "DELETE"},
	{PType: "p", V0: "888", V1: "/sysDictionary/updateSysDictionary", V2: "PUT"},
	{PType: "p", V0: "888", V1: "/sysDictionary/findSysDictionary", V2: "GET"},
	{PType: "p", V0: "888", V1: "/sysDictionary/getSysDictionaryList", V2: "GET"},
	{PType: "p", V0: "888", V1: "/sysOperationRecord/createSysOperationRecord", V2: "POST"},
	{PType: "p", V0: "888", V1: "/sysOperationRecord/deleteSysOperationRecord", V2: "DELETE"},
	{PType: "p", V0: "888", V1: "/sysOperationRecord/updateSysOperationRecord", V2: "PUT"},
	{PType: "p", V0: "888", V1: "/sysOperationRecord/findSysOperationRecord", V2: "GET"},
	{PType: "p", V0: "888", V1: "/sysOperationRecord/getSysOperationRecordList", V2: "GET"},
	{PType: "p", V0: "888", V1: "/sysOperationRecord/deleteSysOperationRecordByIds", V2: "DELETE"},
	{PType: "p", V0: "888", V1: "/user/setUserInfo", V2: "PUT"},
	{PType: "p", V0: "888", V1: "/rule/deleteRule", V2: "DELETE"},
	{PType: "p", V0: "888", V1: "/rule/createRule", V2: "POST"},
	{PType: "p", V0: "888", V1: "/rule/deleteRuleByIds", V2: "DELETE"},
	{PType: "p", V0: "888", V1: "/rule/updateRule", V2: "PUT"},
	{PType: "p", V0: "888", V1: "/rule/findRule", V2: "GET"},
	{PType: "p", V0: "888", V1: "/rule/getRuleList", V2: "GET"},
	{PType: "p", V0: "888", V1: "/rule/switchRuleStatus", V2: "POST"},
	{PType: "p", V0: "888", V1: "/rule/uploadRules", V2: "POST"},
	{PType: "p", V0: "888", V1: "/token/createToken", V2: "POST"},
	{PType: "p", V0: "888", V1: "/token/deleteToken", V2: "DELETE"},
	{PType: "p", V0: "888", V1: "/token/deleteTokenByIds", V2: "DELETE"},
	{PType: "p", V0: "888", V1: "/token/updateToken", V2: "PUT"},
	{PType: "p", V0: "888", V1: "/token/findToken", V2: "GET"},
	{PType: "p", V0: "888", V1: "/token/getTokenList", V2: "GET"},
	{PType: "p", V0: "888", V1: "/searchResult/createSearchResult", V2: "POST"},
	{PType: "p", V0: "888", V1: "/searchResult/deleteSearchResult", V2: "DELETE"},
	{PType: "p", V0: "888", V1: "/searchResult/deleteSearchResultByIds", V2: "DELETE"},
	{PType: "p", V0: "888", V1: "/searchResult/updateSearchResult", V2: "POST"},
	{PType: "p", V0: "888", V1: "/searchResult/findSearchResult", V2: "GET"},
	{PType: "p", V0: "888", V1: "/searchResult/getSearchResultList", V2: "GET"},
	{PType: "p", V0: "888", V1: "/searchResult/exportSearchResult", V2: "GET"},
	{PType: "p", V0: "888", V1: "/searchResult/updateSearchResultStatusByIds", V2: "POST"},
	{PType: "p", V0: "888", V1: "/searchResult/getTaskStatus", V2: "GET"},
	{PType: "p", V0: "888", V1: "/subdomain/createSubdomain", V2: "POST"},
	{PType: "p", V0: "888", V1: "/subdomain/deleteSubdomain", V2: "DELETE"},
	{PType: "p", V0: "888", V1: "/subdomain/deleteSubdomainByIds", V2: "DELETE"},
	{PType: "p", V0: "888", V1: "/subdomain/updateSubdomain", V2: "PUT"},
	{PType: "p", V0: "888", V1: "/subdomain/findSubdomain", V2: "GET"},
	{PType: "p", V0: "888", V1: "/subdomain/getSubdomainList", V2: "GET"},
	{PType: "p", V0: "888", V1: "/filter/createFilter", V2: "POST"},
	{PType: "p", V0: "888", V1: "/filter/deleteFilter", V2: "DELETE"},
	{PType: "p", V0: "888", V1: "/filter/deleteFilterByIds", V2: "DELETE"},
	{PType: "p", V0: "888", V1: "/filter/updateFilter", V2: "PUT"},
	{PType: "p", V0: "888", V1: "/filter/findFilter", V2: "GET"},
	{PType: "p", V0: "888", V1: "/filter/getFilterList", V2: "GET"},
	{PType: "p", V0: "888", V1: "/task/getTaskList", V2: "GET"},
	{PType: "p", V0: "888", V1: "/task/createTask", V2: "POST"},
	{PType: "p", V0: "888", V1: "/task/switchTaskStatus", V2: "POST"},
}

func (c *casbin) Init() error {
	global.GVA_DB.AutoMigrate(gormadapter.CasbinRule{})
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if tx.Find(&[]gormadapter.CasbinRule{}).RowsAffected > 1 {
			color.Danger.Println("\n[Mysql] --> casbin_rule 表的初始数据已存在!")
			return nil
		}
		if err := tx.Create(&carbines).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		color.Info.Println("\n[Mysql] --> casbin_rule 表初始数据成功!")
		return nil
	})
}
