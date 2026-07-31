package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model/request"
	"github.com/madneal/gshark/model/response"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
)

func GetScanLogList(c *gin.Context) {
	var pageInfo request.ScanLogSearch
	if !bindQuery(c, &pageInfo) {
		return
	}
	err, list, total := service.GetScanLogInfoList(pageInfo)
	respondPage(c, err, list, total, pageInfo.Page, pageInfo.PageSize)
}

func GetScanLogOverview(c *gin.Context) {
	overview, err := service.GetScanLogOverview(time.Now())
	if err != nil {
		global.GVA_LOG.Error("get scan log overview failed", zap.Error(err))
		response.FailWithMessage("获取扫描状态失败", c)
		return
	}
	response.OkWithData(overview, c)
}
