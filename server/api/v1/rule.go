package v1

import (
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"github.com/madneal/gshark/model/response"
	"github.com/madneal/gshark/service"
	"go.uber.org/zap"
)

func CreateRule(c *gin.Context) {
	var rule model.Rule
	_ = c.ShouldBindJSON(&rule)
	if err := service.CreateRule(rule); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Any("err", err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

func UploadRules(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("上传失败！", zap.Error(err))
		response.FailWithMessage("上传文件失败", c)
	}
	file, _ := fileHeader.Open()
	csvLines, readErr := csv.NewReader(file).ReadAll()
	if readErr != nil {
		global.GVA_LOG.Error("csv 文件读取失败！", zap.Error(err))
		response.FailWithMessage("规则导入失败", c)
	}
	rules := convertCsvIntoRules(csvLines)
	for _, rule := range rules {
		if err := service.CreateRule(rule); err != nil {
			global.GVA_LOG.Error("创建Rule失败！", zap.Error(err))
			response.FailWithMessage("创建规则失败", c)
			return
		}
	}
}

func convertCsvIntoRules(lines [][]string) []model.Rule {
	rules := make([]model.Rule, 0)
	for index, line := range lines {
		if index == 0 {
			continue
		}
		rules = append(rules, model.Rule{
			RuleType: line[0],
			Content:  line[1],
			Name:     line[2],
			Desc:     line[3],
			Status:   true,
		})
	}
	return rules
}

func DeleteRule(c *gin.Context) {
	var rule model.Rule
	_ = c.ShouldBindJSON(&rule)
	if err := service.DeleteRule(rule); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

func DeleteRuleByIds(c *gin.Context) {
	var IDS request.IdsReq
	_ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteRuleByIds(IDS); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

func UpdateRule(c *gin.Context) {
	var rule model.Rule
	_ = c.ShouldBindJSON(&rule)
	if err := service.UpdateRule(rule); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

func FindRule(c *gin.Context) {
	var rule model.Rule
	_ = c.ShouldBindQuery(&rule)
	if err, rule := service.GetRule(rule.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rule": rule}, c)
	}
}

func GetRuleList(c *gin.Context) {
	var pageInfo request.RuleSearch
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := service.GetRuleInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

func SwitchRuleStatus(c *gin.Context) {
	var switchRequest request.RuleSwitch
	_ = c.ShouldBindJSON(&switchRequest)
	fmt.Println(switchRequest.ID)
	fmt.Println(switchRequest.Status)
	if err := service.SwitchRuleStatus(switchRequest.ID, switchRequest.Status); err != nil {
		global.GVA_LOG.Error("切换状态失败", zap.Any("err", err))
		response.FailWithMessage("切换状态失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}
