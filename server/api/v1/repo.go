package v1

import (
	"github.com/madneal/gshark/global"
    "github.com/madneal/gshark/model"
    "github.com/madneal/gshark/model/request"
    "github.com/madneal/gshark/model/response"
    "github.com/madneal/gshark/service"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

// @Tags Repo
// @Summary 创建Repo
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Repo true "创建Repo"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /repo/createRepo [post]
func CreateRepo(c *gin.Context) {
	var repo model.Repo
	_ = c.ShouldBindJSON(&repo)
	if err := service.CreateRepo(repo); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Any("err", err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// @Tags Repo
// @Summary 删除Repo
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Repo true "删除Repo"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /repo/deleteRepo [delete]
func DeleteRepo(c *gin.Context) {
	var repo model.Repo
	_ = c.ShouldBindJSON(&repo)
	if err := service.DeleteRepo(repo); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Repo
// @Summary 批量删除Repo
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Repo"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /repo/deleteRepoByIds [delete]
func DeleteRepoByIds(c *gin.Context) {
	var IDS request.IdsReq
    _ = c.ShouldBindJSON(&IDS)
	if err := service.DeleteRepoByIds(IDS); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Any("err", err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// @Tags Repo
// @Summary 更新Repo
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Repo true "更新Repo"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /repo/updateRepo [put]
func UpdateRepo(c *gin.Context) {
	var repo model.Repo
	_ = c.ShouldBindJSON(&repo)
	if err := service.UpdateRepo(repo); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Repo
// @Summary 用id查询Repo
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Repo true "用id查询Repo"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /repo/findRepo [get]
func FindRepo(c *gin.Context) {
	var repo model.Repo
	_ = c.ShouldBindQuery(&repo)
	if err, rerepo := service.GetRepo(repo.ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Any("err", err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rerepo": rerepo}, c)
	}
}

// @Tags Repo
// @Summary 分页获取Repo列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.RepoSearch true "分页获取Repo列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /repo/getRepoList [get]
func GetRepoList(c *gin.Context) {
	var pageInfo request.RepoSearch
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := service.GetRepoInfoList(pageInfo); err != nil {
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
