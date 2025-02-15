package initialize

import (
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/middleware"
	"github.com/madneal/gshark/router"
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
		router.InitBaseRouter(PublicGroup) // 注册基础功能路由 不做鉴权
		router.InitInitRouter(PublicGroup) // 自动初始化相关
	}
	PrivateGroup := Router.Group("")
	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		router.InitApiRouter(PrivateGroup)                 // 注册功能api路由
		router.InitJwtRouter(PrivateGroup)                 // jwt相关路由
		router.InitUserRouter(PrivateGroup)                // 注册用户路由
		router.InitMenuRouter(PrivateGroup)                // 注册menu路由
		router.InitSystemRouter(PrivateGroup)              // system相关路由
		router.InitCasbinRouter(PrivateGroup)              // 权限相关路由
		router.InitAuthorityRouter(PrivateGroup)           // 注册角色路由
		router.InitSysDictionaryRouter(PrivateGroup)       // 字典管理
		router.InitSysOperationRecordRouter(PrivateGroup)  // 操作记录
		router.InitSysDictionaryDetailRouter(PrivateGroup) // 字典详情管理
		router.InitRuleRouter(PrivateGroup)
		router.InitTokenRouter(PrivateGroup)
		router.InitSearchResultRouter(PrivateGroup)
		router.InitSubdomainRouter(PrivateGroup)
		router.InitFilterRouter(PrivateGroup)
		router.InitRepoRouter(PrivateGroup)
	}
	global.GVA_LOG.Info("router register success")
	return Router
}
