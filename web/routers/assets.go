/*

Copyright (c) 2018 sec.xiaomi.com

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THEq
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.

*/

package routers

import (
	"github.com/neal1991/x-patrol/models"
	"github.com/neal1991/x-patrol/vars"
	"github.com/go-macaron/session"
	"github.com/go-macaron/csrf"
	"gopkg.in/macaron.v1"
	"github.com/neal1991/x-patrol/util/common"
	"strconv"
	"strings"
)

func ListAssets(ctx *macaron.Context, sess session.Store) {
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
		ctx.HTML(200, "assets")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func NewAssets(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		ctx.HTML(200, "assets_new")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoNewAssets(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		Type := strings.TrimSpace(ctx.Req.Form.Get("type"))
		content := strings.TrimSpace(ctx.Req.Form.Get("content"))
		desc := strings.TrimSpace(ctx.Req.Form.Get("desc"))
		assets := models.NewInputInfo(Type, content, desc)
		assets.Insert()
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
		ctx.HTML(200, "assets_edit")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoEditAssets(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		Type := strings.TrimSpace(ctx.Req.Form.Get("type"))
		content := strings.TrimSpace(ctx.Req.Form.Get("content"))
		desc := strings.TrimSpace(ctx.Req.Form.Get("desc"))
		models.EditInputInfoById(int64(Id), Type, content, desc)
		ctx.Redirect("/admin/assets/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DeleteAssets(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.DeleteInputInfoById(int64(Id))
		ctx.Redirect("/admin/assets/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DeleteAllAssets(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		models.DeleteAllInputInfo()
		ctx.Redirect("/admin/assets/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}
