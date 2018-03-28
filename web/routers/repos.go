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
	"../../models"
	"../../util"

	"gopkg.in/macaron.v1"

	"github.com/go-macaron/session"

	"strconv"
)

func ListRepos(ctx *macaron.Context, sess session.Store) {
	page := ctx.Params(":page")
	p, _ := strconv.Atoi(page)
	p, pre, next := util.GetPreAndNext(p)

	if sess.Get("admin") != nil {
		repos, pages, _ := models.ListReposPage(p)
		pageList := util.GetPageList(p, pages)

		ctx.Data["pages"] = pages
		ctx.Data["page"] = p
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.Data["repos"] = repos
		ctx.HTML(200, "repos")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func EnableRepo(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.EnableRepoById(int64(Id))
		ctx.Redirect("/admin/repos/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DisableRepo(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.DisableRepoById(int64(Id))
		ctx.Redirect("/admin/repos/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DeleteAllRepo(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		models.DeleteAllRepos()
		ctx.Redirect("/admin/repos/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}
