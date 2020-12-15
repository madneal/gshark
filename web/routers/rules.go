package routers

import (
	"github.com/madneal/gshark/models"
	"github.com/madneal/gshark/util"
	"github.com/madneal/gshark/vars"
	"gopkg.in/macaron.v1"

	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"

	"strconv"
	"strings"
)

func ListRules(ctx *macaron.Context, sess session.Store) {
	page := ctx.Params(":page")
	p, _ := strconv.Atoi(page)
	p, pre, next := util.GetPreAndNext(p)

	if sess.Get("admin") != nil {
		rules, pages, _ := models.GetRulesPage(p)
		pageList := util.GetPageList(p, vars.PageStep, pages)

		ctx.Data["pages"] = pages
		ctx.Data["page"] = p
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.Data["rules"] = rules
		ctx.Data["lastPage"] = util.GetLastPage(&pageList)
		ctx.Data["link"] = "/admin/rules/list"
		ctx.Data["role"] = sess.Get("user").(string)
		ctx.HTML(200, "rules")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func NewRules(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		ctx.Data["role"] = sess.Get("user").(string)
		ctx.HTML(200, "rules_new")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoNewRules(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		Type := strings.TrimSpace(ctx.Req.Form.Get("type"))
		content := strings.TrimSpace(ctx.Req.Form.Get("content"))
		caption := strings.TrimSpace(ctx.Req.Form.Get("caption"))
		pos := strings.TrimSpace(ctx.Req.Form.Get("position"))
		desc := strings.TrimSpace(ctx.Req.Form.Get("desc"))
		status := strings.TrimSpace(ctx.Req.Form.Get("status"))
		intStatus, _ := strconv.Atoi(status)
		rule := models.NewRule(Type, content, caption, pos, desc, intStatus)
		rule.Insert()
		ctx.Redirect("/admin/rules/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func EditRules(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		rules, _, _ := models.GetRuleById(int64(Id))
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["rules"] = rules
		ctx.Data["user"] = sess.Get("admin")
		ctx.HTML(200, "rules_edit")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoEditRules(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		position := strings.TrimSpace(ctx.Req.Form.Get("position"))
		Type := strings.TrimSpace(ctx.Req.Form.Get("type"))
		content := strings.TrimSpace(ctx.Req.Form.Get("content"))
		caption := strings.TrimSpace(ctx.Req.Form.Get("caption"))
		desc := strings.TrimSpace(ctx.Req.Form.Get("desc"))
		status := strings.TrimSpace(ctx.Req.Form.Get("status"))
		intStatus, _ := strconv.Atoi(status)
		models.EditRuleById(int64(Id), position, Type, content, caption, desc, intStatus)
		ctx.Redirect("/admin/rules/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DeleteRules(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.DeleteRulesById(int64(Id))
		ctx.Redirect("/admin/rules/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func EnableRules(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.EnableRulesById(int64(Id))
		ctx.Redirect("/admin/rules/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DisableRules(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.DisableRulesById(int64(Id))
		ctx.Redirect("/admin/rules/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}
