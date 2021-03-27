package router

import (
	"gin-vue-admin/api/v1"
	"gin-vue-admin/middleware"
	"github.com/gin-gonic/gin"
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
