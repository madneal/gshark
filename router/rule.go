package router

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/api/v1"
	"github.com/madneal/gshark/middleware"
)

func InitRuleRouter(Router *gin.RouterGroup) {
	RuleRouter := Router.Group("rule").Use(middleware.OperationRecord())
	{
		RuleRouter.POST("createRule", v1.CreateRule)             // 新建Rule
		RuleRouter.DELETE("deleteRule", v1.DeleteRule)           // 删除Rule
		RuleRouter.DELETE("deleteRuleByIds", v1.DeleteRuleByIds) // 批量删除Rule
		RuleRouter.PUT("updateRule", v1.UpdateRule)              // 更新Rule
		RuleRouter.GET("findRule", v1.FindRule)                  // 根据ID获取Rule
		RuleRouter.GET("getRuleList", v1.GetRuleList)            // 获取Rule列表
	}
}
