package source

import (
	"time"

	"github.com/gookit/color"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"

	"gorm.io/gorm"
)

var Api = new(api)

type api struct{}

var apis = []model.SysApi{
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/base/login", "用户登录", "base", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/user/register", "用户注册", "user", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/user/deleteUser", "删除用户", "user", "DELETE"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/user/changePassword", "修改密码", "user", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/user/getUserList", "获取用户列表", "user", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/user/setUserAuthority", "修改用户角色", "user", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/user/setUserInfo", "设置用户信息", "user", "PUT"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/api/createApi", "创建api", "api", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/api/getApiList", "获取api列表", "api", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/api/getApiById", "获取api详细信息", "api", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/api/deleteApi", "删除Api", "api", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/api/updateApi", "更新Api", "api", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/api/getAllApis", "获取所有api", "api", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/authority/createAuthority", "创建角色", "authority", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/authority/deleteAuthority", "删除角色", "authority", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/authority/getAuthorityList", "获取角色列表", "authority", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/getMenu", "获取菜单树", "menu", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/getMenuList", "分页获取基础menu列表", "menu", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/addBaseMenu", "新增菜单", "menu", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/getBaseMenuTree", "获取用户动态路由", "menu", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/addMenuAuthority", "增加menu和角色关联关系", "menu", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/getMenuAuthority", "获取指定角色menu", "menu", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/deleteBaseMenu", "删除菜单", "menu", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/updateBaseMenu", "更新菜单", "menu", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/getBaseMenuById", "根据id获取菜单", "menu", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/casbin/updateCasbin", "更改角色api权限", "casbin", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/casbin/getPolicyPathByAuthorityId", "获取权限列表", "casbin", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/casbin/casbinTest/:pathParam", "RESTFUL模式测试", "casbin", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/jwt/jsonInBlacklist", "jwt加入黑名单(退出)", "jwt", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/authority/setDataAuthority", "设置角色资源权限", "authority", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/authority/updateAuthority", "更新角色信息", "authority", "PUT"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/authority/copyAuthority", "拷贝角色", "authority", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/system/getSystemConfig", "获取配置文件内容", "system", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/system/setSystemConfig", "设置配置文件内容", "system", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/system/getServerInfo", "获取服务器信息", "system", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/system/emailTest", "发送测试邮件", "system", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/system/botTest", "发送bot测试消息", "system", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionaryDetail/createSysDictionaryDetail", "新增字典内容", "sysDictionaryDetail", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionaryDetail/deleteSysDictionaryDetail", "删除字典内容", "sysDictionaryDetail", "DELETE"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionaryDetail/updateSysDictionaryDetail", "更新字典内容", "sysDictionaryDetail", "PUT"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionaryDetail/findSysDictionaryDetail", "根据ID获取字典内容", "sysDictionaryDetail", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionaryDetail/getSysDictionaryDetailList", "获取字典内容列表", "sysDictionaryDetail", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionary/createSysDictionary", "新增字典", "sysDictionary", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionary/deleteSysDictionary", "删除字典", "sysDictionary", "DELETE"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionary/updateSysDictionary", "更新字典", "sysDictionary", "PUT"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionary/findSysDictionary", "根据ID获取字典", "sysDictionary", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionary/getSysDictionaryList", "获取字典列表", "sysDictionary", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysOperationRecord/createSysOperationRecord", "新增操作记录", "sysOperationRecord", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysOperationRecord/deleteSysOperationRecord", "删除操作记录", "sysOperationRecord", "DELETE"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysOperationRecord/findSysOperationRecord", "根据ID获取操作记录", "sysOperationRecord", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysOperationRecord/getSysOperationRecordList", "获取操作记录列表", "sysOperationRecord", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysOperationRecord/deleteSysOperationRecordByIds", "批量删除操作历史", "sysOperationRecord", "DELETE"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/rule/deleteRule", "删除规则", "rule", "DELETE"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/rule/createRule", "新增规则", "rule", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/rule/deleteRuleByIds", "批量删除规则", "rule", "DELETE"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/rule/updateRule", "更新规则", "rule", "PUT"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/rule/findRule", "根据ID获取规则", "rule", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/rule/getRuleList", "获取规则列表", "rule", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/rule/uploadRules", "规则导入", "rule", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/rule/switchRuleStatus", "变更规则状态", "rule", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/token/createToken", "新增token", "token", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/token/deleteToken", "删除token", "token", "DELETE"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/token/deleteTokenByIds", "批量删除token", "token", "DELETE"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/token/updateToken", "更新token", "token", "PUT"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/token/findToken", "根据ID获取token", "token", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/token/getTokenList", "获取token列表", "token", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/searchResult/createSearchResult", "新增搜索结果", "searchResult", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/searchResult/deleteSearchResult", "删除搜索结果", "searchResult", "DELETE"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/searchResult/deleteSearchResultByIds", "批量删除搜索结果", "searchResult", "DELETE"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/searchResult/updateSearchResult", "更新搜索结果", "searchResult", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/searchResult/findSearchResult", "根据ID获取搜索结果", "searchResult", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/searchResult/getSearchResultList", "获取搜索结果列表", "searchResult", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/searchResult/exportSearchResult", "导出结果列表", "searchResult", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/searchResult/updateSearchResultStatusByIds", "批量更新搜索结果列表", "searchResult", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/searchResult/startSecFilterTask", "开始二次过滤任务", "searchResult", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/searchResult/getTaskStatus", "获取任务状态", "searchResult", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/subdomain/createSubdomain", "新增子域名", "subdomain", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/subdomain/deleteSubdomain", "删除子域名", "subdomain", "DELETE"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/subdomain/deleteSubdomainByIds", "批量删除子域名", "subdomain", "DELETE"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/subdomain/updateSubdomain", "更新子域名", "subdomain", "PUT"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/subdomain/findSubdomain", "根据ID获取子域名", "subdomain", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/subdomain/getSubdomainList", "获取子域名列表", "subdomain", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/filter/createFilter", "新增过滤规则", "filter", "POST"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/filter/deleteFilter", "删除过滤规则", "filter", "DELETE"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/filter/deleteFilterByIds", "批量删除过滤规则", "filter", "DELETE"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/filter/updateFilter", "更新过滤规则", "filter", "PUT"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/filter/findFilter", "根据ID获取过滤规则", "filter", "GET"},
	{global.GVA_MODEL{CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/filter/getFilterList", "获取过滤规则列表", "filter", "GET"},
}

func (a *api) Init() error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if tx.Where("id IN ?", []int{1, 67}).Find(&[]model.SysApi{}).RowsAffected == 2 {
			color.Danger.Println("\n[Mysql] --> sys_apis 表的初始数据已存在!")
			return nil
		}
		if err := tx.Create(&apis).Error; err != nil { // 遇到错误时回滚事务
			return err
		}
		color.Info.Println("\n[Mysql] --> sys_apis 表初始数据成功!")
		return nil
	})
}
