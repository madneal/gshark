package routers

import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/session"
	"strconv"
	"github.com/neal1991/gshark/util/common"
	"github.com/neal1991/gshark/models"

	"github.com/neal1991/gshark/vars"
)

func ListAppAssets(ctx *macaron.Context, sess session.Store) {
	page := ctx.Params(":page")
	p, _ := strconv.Atoi(page)
	p, pre, next := common.GetPreAndNext(p)

	if sess.Get("admin") != nil {
		assets, pages, _ := models.ListInputInfoPage(p)
		pageList := common.GetPageList(p, vars.PageStep, pages)

		ctx.Data["pages"] = pages
		ctx.Data["page"] = p
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.Data["assets"] = assets
		ctx.HTML(200, "app_assets")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DetectApp(ctx *macaron.Context, sess session.Store)  {
	if sess.Get("admin") != nil {
		hash := ctx.Query("hash")
		isExist, id := models.Detect(hash)
		ctx.JSON(200, map[string]interface{} {
			"isExist": isExist,
			"id": id,
		})
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func GetAppAsset(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id_ := ctx.Query("id")
		id, _ := strconv.ParseInt(id_, 10, 64)
		appAsset := models.GetAppAssetById(id)
		ctx.Data["appAsset"] = appAsset
		ctx.HTML(200, "app_asset_detail")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

