package source

import (
	"github.com/gookit/color"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"time"

	"gorm.io/gorm"
)

var Api = new(api)

type api struct{}

var apis = []model.SysApi{
	{global.GVA_MODEL{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/base/login", "用户登录", "base", "POST"},
	{global.GVA_MODEL{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/user/register", "用户注册", "user", "POST"},
	{global.GVA_MODEL{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/api/createApi", "创建api", "api", "POST"},
	{global.GVA_MODEL{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/api/getApiList", "获取api列表", "api", "POST"},
	{global.GVA_MODEL{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/api/getApiById", "获取api详细信息", "api", "POST"},
	{global.GVA_MODEL{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/api/deleteApi", "删除Api", "api", "POST"},
	{global.GVA_MODEL{ID: 7, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/api/updateApi", "更新Api", "api", "POST"},
	{global.GVA_MODEL{ID: 8, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/api/getAllApis", "获取所有api", "api", "POST"},
	{global.GVA_MODEL{ID: 9, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/authority/createAuthority", "创建角色", "authority", "POST"},
	{global.GVA_MODEL{ID: 10, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/authority/deleteAuthority", "删除角色", "authority", "POST"},
	{global.GVA_MODEL{ID: 11, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/authority/getAuthorityList", "获取角色列表", "authority", "POST"},
	{global.GVA_MODEL{ID: 12, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/getMenu", "获取菜单树", "menu", "POST"},
	{global.GVA_MODEL{ID: 13, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/getMenuList", "分页获取基础menu列表", "menu", "POST"},
	{global.GVA_MODEL{ID: 14, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/addBaseMenu", "新增菜单", "menu", "POST"},
	{global.GVA_MODEL{ID: 15, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/getBaseMenuTree", "获取用户动态路由", "menu", "POST"},
	{global.GVA_MODEL{ID: 16, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/addMenuAuthority", "增加menu和角色关联关系", "menu", "POST"},
	{global.GVA_MODEL{ID: 17, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/getMenuAuthority", "获取指定角色menu", "menu", "POST"},
	{global.GVA_MODEL{ID: 18, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/deleteBaseMenu", "删除菜单", "menu", "POST"},
	{global.GVA_MODEL{ID: 19, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/updateBaseMenu", "更新菜单", "menu", "POST"},
	{global.GVA_MODEL{ID: 20, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/menu/getBaseMenuById", "根据id获取菜单", "menu", "POST"},
	{global.GVA_MODEL{ID: 21, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/user/changePassword", "修改密码", "user", "POST"},
	{global.GVA_MODEL{ID: 23, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/user/getUserList", "获取用户列表", "user", "POST"},
	{global.GVA_MODEL{ID: 24, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/user/setUserAuthority", "修改用户角色", "user", "POST"},
	{global.GVA_MODEL{ID: 25, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/fileUploadAndDownload/upload", "文件上传示例", "fileUploadAndDownload", "POST"},
	{global.GVA_MODEL{ID: 26, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/fileUploadAndDownload/getFileList", "获取上传文件列表", "fileUploadAndDownload", "POST"},
	{global.GVA_MODEL{ID: 27, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/casbin/updateCasbin", "更改角色api权限", "casbin", "POST"},
	{global.GVA_MODEL{ID: 28, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/casbin/getPolicyPathByAuthorityId", "获取权限列表", "casbin", "POST"},
	{global.GVA_MODEL{ID: 29, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/fileUploadAndDownload/deleteFile", "删除文件", "fileUploadAndDownload", "POST"},
	{global.GVA_MODEL{ID: 30, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/jwt/jsonInBlacklist", "jwt加入黑名单(退出)", "jwt", "POST"},
	{global.GVA_MODEL{ID: 31, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/authority/setDataAuthority", "设置角色资源权限", "authority", "POST"},
	{global.GVA_MODEL{ID: 32, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/system/getSystemConfig", "获取配置文件内容", "system", "POST"},
	{global.GVA_MODEL{ID: 33, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/system/setSystemConfig", "设置配置文件内容", "system", "POST"},
	{global.GVA_MODEL{ID: 39, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/casbin/casbinTest/:pathParam", "RESTFUL模式测试", "casbin", "GET"},
	{global.GVA_MODEL{ID: 40, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/autoCode/createTemp", "自动化代码", "autoCode", "POST"},
	{global.GVA_MODEL{ID: 41, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/authority/updateAuthority", "更新角色信息", "authority", "PUT"},
	{global.GVA_MODEL{ID: 42, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/authority/copyAuthority", "拷贝角色", "authority", "POST"},
	{global.GVA_MODEL{ID: 43, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/user/deleteUser", "删除用户", "user", "DELETE"},
	{global.GVA_MODEL{ID: 44, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionaryDetail/createSysDictionaryDetail", "新增字典内容", "sysDictionaryDetail", "POST"},
	{global.GVA_MODEL{ID: 45, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionaryDetail/deleteSysDictionaryDetail", "删除字典内容", "sysDictionaryDetail", "DELETE"},
	{global.GVA_MODEL{ID: 46, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionaryDetail/updateSysDictionaryDetail", "更新字典内容", "sysDictionaryDetail", "PUT"},
	{global.GVA_MODEL{ID: 47, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionaryDetail/findSysDictionaryDetail", "根据ID获取字典内容", "sysDictionaryDetail", "GET"},
	{global.GVA_MODEL{ID: 48, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionaryDetail/getSysDictionaryDetailList", "获取字典内容列表", "sysDictionaryDetail", "GET"},
	{global.GVA_MODEL{ID: 49, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionary/createSysDictionary", "新增字典", "sysDictionary", "POST"},
	{global.GVA_MODEL{ID: 50, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionary/deleteSysDictionary", "删除字典", "sysDictionary", "DELETE"},
	{global.GVA_MODEL{ID: 51, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionary/updateSysDictionary", "更新字典", "sysDictionary", "PUT"},
	{global.GVA_MODEL{ID: 52, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionary/findSysDictionary", "根据ID获取字典", "sysDictionary", "GET"},
	{global.GVA_MODEL{ID: 53, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysDictionary/getSysDictionaryList", "获取字典列表", "sysDictionary", "GET"},
	{global.GVA_MODEL{ID: 54, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysOperationRecord/createSysOperationRecord", "新增操作记录", "sysOperationRecord", "POST"},
	{global.GVA_MODEL{ID: 55, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysOperationRecord/deleteSysOperationRecord", "删除操作记录", "sysOperationRecord", "DELETE"},
	{global.GVA_MODEL{ID: 56, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysOperationRecord/findSysOperationRecord", "根据ID获取操作记录", "sysOperationRecord", "GET"},
	{global.GVA_MODEL{ID: 57, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysOperationRecord/getSysOperationRecordList", "获取操作记录列表", "sysOperationRecord", "GET"},
	{global.GVA_MODEL{ID: 58, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/autoCode/getTables", "获取数据库表", "autoCode", "GET"},
	{global.GVA_MODEL{ID: 59, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/autoCode/getDB", "获取所有数据库", "autoCode", "GET"},
	{global.GVA_MODEL{ID: 60, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/autoCode/getColumn", "获取所选table的所有字段", "autoCode", "GET"},
	{global.GVA_MODEL{ID: 61, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/sysOperationRecord/deleteSysOperationRecordByIds", "批量删除操作历史", "sysOperationRecord", "DELETE"},
	{global.GVA_MODEL{ID: 62, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/simpleUploader/upload", "插件版分片上传", "simpleUploader", "POST"},
	{global.GVA_MODEL{ID: 63, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/simpleUploader/checkFileMd5", "文件完整度验证", "simpleUploader", "GET"},
	{global.GVA_MODEL{ID: 64, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/simpleUploader/mergeFileMd5", "上传完成合并文件", "simpleUploader", "GET"},
	{global.GVA_MODEL{ID: 65, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/user/setUserInfo", "设置用户信息", "user", "PUT"},
	{global.GVA_MODEL{ID: 66, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/system/getServerInfo", "获取服务器信息", "system", "POST"},
	{global.GVA_MODEL{ID: 80, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/autoCode/preview", "预览自动化代码", "autoCode", "POST"},
	{global.GVA_MODEL{ID: 81, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/excel/importExcel", "导入excel", "excel", "POST"},
	{global.GVA_MODEL{ID: 82, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/excel/loadExcel", "下载excel", "excel", "GET"},
	{global.GVA_MODEL{ID: 83, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/excel/exportExcel", "导出excel", "excel", "POST"},
	{global.GVA_MODEL{ID: 84, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/excel/downloadTemplate", "下载excel模板", "excel", "GET"},
	{global.GVA_MODEL{ID: 85, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/excel/downloadTemplate", "下载excel模板", "excel", "GET"},
	{global.GVA_MODEL{ID: 86, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/rule/deleteRule", "删除规则", "rule", "DELETE"},
	{global.GVA_MODEL{ID: 87, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/rule/deleteRuleByIds", "批量删除规则", "rule", "DELETE"},
	{global.GVA_MODEL{ID: 88, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/rule/updateRule", "更新规则", "rule", "PUT"},
	{global.GVA_MODEL{ID: 89, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/rule/findRule", "根据ID获取规则", "rule", "GET"},
	{global.GVA_MODEL{ID: 90, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/rule/getRuleList", "获取规则列表", "rule", "GET"},
	{global.GVA_MODEL{ID: 91, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/token/createToken", "新增token", "token", "POST"},
	{global.GVA_MODEL{ID: 92, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/token/deleteToken", "删除token", "token", "DELETE"},
	{global.GVA_MODEL{ID: 93, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/token/deleteTokenByIds", "批量删除token", "token", "DELETE"},
	{global.GVA_MODEL{ID: 94, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/token/updateToken", "更新token", "token", "PUT"},
	{global.GVA_MODEL{ID: 95, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/token/findToken", "根据ID获取token", "token", "GET"},
	{global.GVA_MODEL{ID: 96, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/token/getTokenList", "获取token列表", "token", "GET"},
	{global.GVA_MODEL{ID: 97, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/searchResult/createSearchResult", "新增搜索结果", "searchResult", "POST"},
	{global.GVA_MODEL{ID: 98, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/searchResult/deleteSearchResult", "删除搜索结果", "searchResult", "DELETE"},
	{global.GVA_MODEL{ID: 99, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/searchResult/deleteSearchResultByIds", "批量删除搜索结果", "searchResult", "DELETE"},
	{global.GVA_MODEL{ID: 100, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/searchResult/updateSearchResult", "更新搜索结果", "searchResult", "POST"},
	{global.GVA_MODEL{ID: 101, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/searchResult/findSearchResult", "根据ID获取搜索结果", "searchResult", "GET"},
	{global.GVA_MODEL{ID: 102, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/searchResult/getSearchResultList", "获取搜索结果列表", "searchResult", "GET"},
	{global.GVA_MODEL{ID: 103, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/subdomain/createSubdomain", "新增子域名", "subdomain", "POST"},
	{global.GVA_MODEL{ID: 104, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/subdomain/deleteSubdomain", "删除子域名", "subdomain", "DELETE"},
	{global.GVA_MODEL{ID: 105, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/subdomain/deleteSubdomainByIds", "批量删除子域名", "subdomain", "DELETE"},
	{global.GVA_MODEL{ID: 106, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/subdomain/updateSubdomain", "更新子域名", "subdomain", "PUT"},
	{global.GVA_MODEL{ID: 107, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/subdomain/findSubdomain", "根据ID获取子域名", "subdomain", "GET"},
	{global.GVA_MODEL{ID: 108, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/subdomain/getSubdomainList", "获取子域名列表", "subdomain", "GET"},
	{global.GVA_MODEL{ID: 115, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/filter/createFilter", "新增过滤规则", "filter", "POST"},
	{global.GVA_MODEL{ID: 116, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/filter/deleteFilter", "删除过滤规则", "filter", "DELETE"},
	{global.GVA_MODEL{ID: 117, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/filter/deleteFilterByIds", "批量删除过滤规则", "filter", "DELETE"},
	{global.GVA_MODEL{ID: 118, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/filter/updateFilter", "更新过滤规则", "filter", "PUT"},
	{global.GVA_MODEL{ID: 119, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/filter/findFilter", "根据ID获取过滤规则", "filter", "GET"},
	{global.GVA_MODEL{ID: 120, CreatedAt: time.Now(), UpdatedAt: time.Now()}, "/filter/getFilterList", "获取过滤规则列表", "filter", "GET"},
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@ sys_apis 表数据初始化
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
