package routers

import (
	"github.com/go-macaron/session"
	"github.com/neal1991/gshark/models"
	"github.com/neal1991/gshark/util/common"
	"gopkg.in/macaron.v1"
	"strconv"

	"fmt"
	"github.com/go-macaron/csrf"
	"github.com/neal1991/gshark/vars"
)

func ListAppAssets(ctx *macaron.Context, sess session.Store) {
	page := ctx.Params(":page")
	p, _ := strconv.Atoi(page)
	p, pre, next := common.GetPreAndNext(p)

	if sess.Get("admin") != nil {
		assets, pages, _ := models.ListAppAssets(p)
		pageList := common.GetPageList(p, vars.PageStep, pages)

		ctx.Data["pages"] = pages
		ctx.Data["page"] = p
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.Data["appAssets"] = assets
		ctx.Data["role"] = sess.Get("user").(string)
		ctx.HTML(200, "app_assets")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DetectApp(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		hash := ctx.Query("hash")
		isExist, id := models.Detect(hash)
		ctx.JSON(200, map[string]interface{}{
			"isExist": isExist,
			"id":      id,
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

func NewAppAsset(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		ctx.Data["role"] = sess.Get("user").(string)
		ctx.HTML(200, "app_asset_new")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoNewAppAsset(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		status, _ := strconv.Atoi(ctx.Query("status"))
		appAsset := models.NewAppAsset(
			ctx.Query("name"),
			ctx.Query("description"),
			ctx.Query("market"),
			ctx.Query("developer"),
			ctx.Query("version"),
			ctx.Query("deployDate"),
			ctx.Query("url"),
			ctx.Query("hash"),
			ctx.Query("filename"),
			status,
		)
		appAsset.Insert()
		ctx.Redirect("/admin/app/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DelAppAsset(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id, _ := strconv.Atoi(ctx.Query("id"))
		models.DeleteAppAssetById(int64(id))
		ctx.Redirect("/admin/app/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func EditAppAsset(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("admin") != nil {
		id, _ := strconv.Atoi(ctx.Query("id"))
		appAsset := models.GetAppAssetById(int64(id))
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["appAsset"] = appAsset
		ctx.Data["role"] = sess.Get("user").(string)
		ctx.HTML(200, "app_asset_edit")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoEditAppAsset(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id, _ := strconv.Atoi(ctx.Query("id"))
		status, _ := strconv.Atoi(ctx.Query("status"))
		appAsset := models.NewAppAsset(
			ctx.Query("name"),
			ctx.Query("description"),
			ctx.Query("market"),
			ctx.Query("developer"),
			ctx.Query("version"),
			ctx.Query("deployDate"),
			ctx.Query("url"),
			ctx.Query("hash"),
			ctx.Query("filename"),
			status,
		)
		err := models.EditAppAssetById(int64(id), appAsset)
		if err != nil {
			fmt.Println(err)
		}
		ctx.Redirect("/admin/app/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}
