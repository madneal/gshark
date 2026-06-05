package source

import (
	"github.com/gookit/color"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"gorm.io/gorm"
)

var Casbin = new(casbin)

type casbin struct{}

var carbines = []model.CasbinModel{
	{Ptype: "p", AuthorityId: "888", Path: "/base/login", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/user/register", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/api/createApi", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/api/getApiList", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/api/getApiById", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/api/deleteApi", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/api/updateApi", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/api/getAllApis", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/authority/createAuthority", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/authority/deleteAuthority", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/authority/getAuthorityList", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/authority/setDataAuthority", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/authority/updateAuthority", Method: "PUT"},
	{Ptype: "p", AuthorityId: "888", Path: "/authority/copyAuthority", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/menu/getMenu", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/menu/getMenuList", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/menu/addBaseMenu", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/menu/getBaseMenuTree", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/menu/addMenuAuthority", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/menu/getMenuAuthority", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/menu/deleteBaseMenu", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/menu/updateBaseMenu", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/menu/getBaseMenuById", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/user/changePassword", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/user/getUserList", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/user/setUserAuthority", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/user/deleteUser", Method: "DELETE"},
	{Ptype: "p", AuthorityId: "888", Path: "/casbin/updateCasbin", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/casbin/getPolicyPathByAuthorityId", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/casbin/casbinTest/:pathParam", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/jwt/jsonInBlacklist", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/system/getSystemConfig", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/system/setSystemConfig", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/system/getServerInfo", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/system/emailTest", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/system/botTest", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/sysDictionaryDetail/createSysDictionaryDetail", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/sysDictionaryDetail/deleteSysDictionaryDetail", Method: "DELETE"},
	{Ptype: "p", AuthorityId: "888", Path: "/sysDictionaryDetail/updateSysDictionaryDetail", Method: "PUT"},
	{Ptype: "p", AuthorityId: "888", Path: "/sysDictionaryDetail/findSysDictionaryDetail", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/sysDictionaryDetail/getSysDictionaryDetailList", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/sysDictionary/createSysDictionary", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/sysDictionary/deleteSysDictionary", Method: "DELETE"},
	{Ptype: "p", AuthorityId: "888", Path: "/sysDictionary/updateSysDictionary", Method: "PUT"},
	{Ptype: "p", AuthorityId: "888", Path: "/sysDictionary/findSysDictionary", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/sysDictionary/getSysDictionaryList", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/sysOperationRecord/createSysOperationRecord", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/sysOperationRecord/deleteSysOperationRecord", Method: "DELETE"},
	{Ptype: "p", AuthorityId: "888", Path: "/sysOperationRecord/updateSysOperationRecord", Method: "PUT"},
	{Ptype: "p", AuthorityId: "888", Path: "/sysOperationRecord/findSysOperationRecord", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/sysOperationRecord/getSysOperationRecordList", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/sysOperationRecord/deleteSysOperationRecordByIds", Method: "DELETE"},
	{Ptype: "p", AuthorityId: "888", Path: "/user/setUserInfo", Method: "PUT"},
	{Ptype: "p", AuthorityId: "888", Path: "/rule/deleteRule", Method: "DELETE"},
	{Ptype: "p", AuthorityId: "888", Path: "/rule/createRule", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/rule/deleteRuleByIds", Method: "DELETE"},
	{Ptype: "p", AuthorityId: "888", Path: "/rule/updateRule", Method: "PUT"},
	{Ptype: "p", AuthorityId: "888", Path: "/rule/findRule", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/rule/getRuleList", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/rule/switchRuleStatus", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/rule/uploadRules", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/token/createToken", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/token/deleteToken", Method: "DELETE"},
	{Ptype: "p", AuthorityId: "888", Path: "/token/deleteTokenByIds", Method: "DELETE"},
	{Ptype: "p", AuthorityId: "888", Path: "/token/updateToken", Method: "PUT"},
	{Ptype: "p", AuthorityId: "888", Path: "/token/findToken", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/token/getTokenList", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/searchResult/createSearchResult", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/searchResult/deleteSearchResult", Method: "DELETE"},
	{Ptype: "p", AuthorityId: "888", Path: "/searchResult/deleteSearchResultByIds", Method: "DELETE"},
	{Ptype: "p", AuthorityId: "888", Path: "/searchResult/updateSearchResult", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/searchResult/findSearchResult", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/searchResult/getSearchResultList", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/searchResult/exportSearchResult", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/searchResult/updateSearchResultStatusByIds", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/searchResult/getTaskStatus", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/subdomain/createSubdomain", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/subdomain/deleteSubdomain", Method: "DELETE"},
	{Ptype: "p", AuthorityId: "888", Path: "/subdomain/deleteSubdomainByIds", Method: "DELETE"},
	{Ptype: "p", AuthorityId: "888", Path: "/subdomain/updateSubdomain", Method: "PUT"},
	{Ptype: "p", AuthorityId: "888", Path: "/subdomain/findSubdomain", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/subdomain/getSubdomainList", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/filter/createFilter", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/filter/deleteFilter", Method: "DELETE"},
	{Ptype: "p", AuthorityId: "888", Path: "/filter/deleteFilterByIds", Method: "DELETE"},
	{Ptype: "p", AuthorityId: "888", Path: "/filter/updateFilter", Method: "PUT"},
	{Ptype: "p", AuthorityId: "888", Path: "/filter/findFilter", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/filter/getFilterList", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/task/getTaskList", Method: "GET"},
	{Ptype: "p", AuthorityId: "888", Path: "/task/createTask", Method: "POST"},
	{Ptype: "p", AuthorityId: "888", Path: "/task/switchTaskStatus", Method: "POST"},
}

func (c *casbin) Init() error {
	global.GVA_DB.AutoMigrate(model.CasbinModel{})
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if tx.Find(&[]model.CasbinModel{}).RowsAffected > 1 {
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
