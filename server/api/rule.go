package api

import (
	"encoding/csv"
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
	if !bindJSON(c, &rule) {
		return
	}
	respondMutation(c, service.CreateRule(rule), "创建失败!", "创建失败", "创建成功")
}

func UploadRules(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		global.GVA_LOG.Error("上传失败！", zap.Error(err))
		response.FailWithMessage("上传文件失败", c)
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		global.GVA_LOG.Error("打开文件失败！", zap.Error(err))
		response.FailWithMessage("打开文件失败", c)
		return
	}
	defer file.Close()
	csvLines, readErr := csv.NewReader(file).ReadAll()
	if readErr != nil {
		global.GVA_LOG.Error("csv 文件读取失败！", zap.Error(readErr))
		response.FailWithMessage("规则导入失败", c)
		return
	}
	rules := convertCsvIntoRules(csvLines)
	for _, rule := range rules {
		if err := service.CreateRule(rule); err != nil {
			global.GVA_LOG.Error("创建规则失败！", zap.Error(err))
			response.FailWithMessage("创建规则失败", c)
			return
		}
	}
	response.OkWithMessage("规则导入成功", c)
}

func convertCsvIntoRules(lines [][]string) []model.Rule {
	rules := make([]model.Rule, 0)
	for index, line := range lines {
		if index == 0 {
			continue
		}
		if len(line) < 4 {
			global.GVA_LOG.Warn("skip invalid rule csv row", zap.Int("row", index+1))
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
	if !bindJSON(c, &rule) {
		return
	}
	respondMutation(c, service.DeleteRule(rule), "删除失败!", "删除失败", "删除成功")
}

func DeleteRuleByIds(c *gin.Context) {
	var IDS request.IdsReq
	if !bindJSON(c, &IDS) {
		return
	}
	respondMutation(c, service.DeleteRuleByIds(IDS), "批量删除失败!", "批量删除失败", "批量删除成功")
}

func UpdateRule(c *gin.Context) {
	var rule model.Rule
	if !bindJSON(c, &rule) {
		return
	}
	respondMutation(c, service.UpdateRule(rule), "更新失败!", "更新失败", "更新成功")
}

func FindRule(c *gin.Context) {
	var rule model.Rule
	if !bindQuery(c, &rule) {
		return
	}
	if err, rule := service.GetRule(rule.ID); err != nil {
		respondMutation(c, err, "查询失败!", "查询失败", "")
	} else {
		response.OkWithData(gin.H{"rule": rule}, c)
	}
}

func GetRuleList(c *gin.Context) {
	var pageInfo request.RuleSearch
	if !bindQuery(c, &pageInfo) {
		return
	}
	err, list, total := service.GetRuleInfoList(pageInfo)
	respondPage(c, err, list, total, pageInfo.Page, pageInfo.PageSize)
}

func SwitchRuleStatus(c *gin.Context) {
	var switchRequest request.RuleSwitch
	if !bindJSON(c, &switchRequest) {
		return
	}
	respondMutation(c, service.SwitchRuleStatus(switchRequest.ID, switchRequest.Status), "切换状态失败", "切换状态失败", "更新成功")
}
