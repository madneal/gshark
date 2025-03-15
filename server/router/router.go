package router

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 初始化总路由

func Routers() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	var Router = gin.Default()
	Router.StaticFS(global.GVA_CONFIG.Local.Path, http.Dir(global.GVA_CONFIG.Local.Path)) // 为用户头像和文件提供静态地址
	// Router.Use(middleware.LoadTls())  // 打开就能玩https了
	global.GVA_LOG.Info("use middleware logger")
	// 跨域
	//Router.Use(middleware.Cors()) // 如需跨域可以打开
	// 方便统一添加路由组前缀 多服务器上线使用
	PublicGroup := Router.Group("")
	{
		InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
		InitInitRouter(PublicGroup) // 自动初始化相关
	}
	PrivateGroup := Router.Group("")
	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		InitApiRouter(PrivateGroup)                 // 注册功能api路由
		InitJwtRouter(PrivateGroup)                 // jwt相关路由
		InitUserRouter(PrivateGroup)                // 注册用户路由
		InitMenuRouter(PrivateGroup)                // 注册menu路由
		InitSystemRouter(PrivateGroup)              // system相关路由
		InitCasbinRouter(PrivateGroup)              // 权限相关路由
		InitAuthorityRouter(PrivateGroup)           // 注册角色路由
		InitSysDictionaryRouter(PrivateGroup)       // 字典管理
		InitSysOperationRecordRouter(PrivateGroup)  // 操作记录
		InitSysDictionaryDetailRouter(PrivateGroup) // 字典详情管理
		InitRuleRouter(PrivateGroup)
		InitTokenRouter(PrivateGroup)
		InitSearchResultRouter(PrivateGroup)
		InitSubdomainRouter(PrivateGroup)
		InitFilterRouter(PrivateGroup)
		InitRepoRouter(PrivateGroup)
	}
	global.GVA_LOG.Info("router register success")
	return Router
}
