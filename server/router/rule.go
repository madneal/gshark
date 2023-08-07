package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api"
	"github.com/madneal/gshark/middleware"
)

func InitRuleRouter(Router *gin.RouterGroup) {
	RuleRouter := Router.Group("rule").Use(middleware.OperationRecord())
	{
		RuleRouter.POST("createRule", api.CreateRule)             // 新建Rule
		RuleRouter.DELETE("deleteRule", api.DeleteRule)           // 删除Rule
		RuleRouter.DELETE("deleteRuleByIds", api.DeleteRuleByIds) // 批量删除Rule
		RuleRouter.PUT("updateRule", api.UpdateRule)              // 更新Rule
		RuleRouter.GET("findRule", api.FindRule)                  // 根据ID获取Rule
		RuleRouter.GET("getRuleList", api.GetRuleList)            // 获取Rule列表
		RuleRouter.POST("switchRuleStatus", api.SwitchRuleStatus)
		RuleRouter.POST("uploadRules", api.UploadRules)
	}
}
