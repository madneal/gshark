package routers

import (
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/madneal/gshark/logger"
	"github.com/madneal/gshark/models"
	"github.com/madneal/gshark/util"
	"github.com/madneal/gshark/vars"
	"gopkg.in/macaron.v1"
	"strconv"
	"strings"
)

func ListAssets(ctx *macaron.Context, sess session.Store) {
	page := ctx.Params(":page")
	p, _ := strconv.Atoi(page)
	p, pre, next := util.GetPreAndNext(p)

	if sess.Get("admin") != nil {
		assets, pages, _ := models.ListInputInfoPage(p)
		pageList := util.GetPageList(p, vars.PageStep, pages)

		ctx.Data["pages"] = pages
		ctx.Data["page"] = p
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.Data["assets"] = assets
		ctx.Data["role"] = sess.Get("user").(string)
		ctx.Data["lastPage"] = util.GetLastPage(&pageList)
		ctx.Data["link"] = "/admin/assets/list"
		ctx.HTML(200, "assets")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func NewAssets(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		ctx.Data["role"] = sess.Get("user").(string)
		ctx.HTML(200, "assets_new")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoNewAssets(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		err := ctx.Req.ParseForm()
		Type := strings.TrimSpace(ctx.Req.Form.Get("type"))
		content := strings.TrimSpace(ctx.Req.Form.Get("content"))
		desc := strings.TrimSpace(ctx.Req.Form.Get("desc"))
		assets := models.NewInputInfo(Type, content, desc)
		_, err = assets.Insert()
		if err != nil {
			logger.Log.Error(err)
		}
		ctx.Redirect("/admin/assets/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func EditAssets(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		assets, _, _ := models.GetInputInfoById(int64(Id))
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["assets"] = assets
		ctx.Data["user"] = sess.Get("admin")
		ctx.Data["role"] = sess.Get("user").(string)
		ctx.HTML(200, "assets_edit")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoEditAssets(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		err := ctx.Req.ParseForm()
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		Type := strings.TrimSpace(ctx.Req.Form.Get("type"))
		content := strings.TrimSpace(ctx.Req.Form.Get("content"))
		desc := strings.TrimSpace(ctx.Req.Form.Get("desc"))
		err = models.EditInputInfoById(int64(Id), Type, content, desc)
		if err != nil {
			logger.Log.Error(err)
		}
		ctx.Redirect("/admin/assets/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DeleteAssets(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		err := models.DeleteInputInfoById(int64(Id))
		if err != nil {
			logger.Log.Error(err)
		}
		ctx.Redirect("/admin/assets/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DeleteAllAssets(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		err := models.DeleteAllInputInfo()
		if err != nil {
			logger.Log.Error(err)
		}
		ctx.Redirect("/admin/assets/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}
