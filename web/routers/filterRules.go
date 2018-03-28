package routers

import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/session"
	"strconv"
	"../../models"
	"../../util"
	"strings"
	"github.com/go-macaron/csrf"
)

func ListFilterRules(ctx *macaron.Context, sess session.Store) {
	page := ctx.Params(":page")
	p, _ := strconv.Atoi(page)
	p, pre, next := util.GetPreAndNext(p)

	if sess.Get("admin") != nil {
		rules, pages, _ := models.GetFilterRulesPage(p)
		pageList := util.GetPageList(p, pages)

		ctx.Data["pages"] = pages
		ctx.Data["page"] = page
		ctx.Data["pre"] = pre
		ctx.Data["next"] = next
		ctx.Data["pageList"] = pageList
		ctx.Data["rules"] = rules
		ctx.HTML(200, "filterrules")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func NewRule(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		ctx.Data["type"] = "new"
		rule := models.NewFilterRule(0, "", "")
		ctx.Data["rule"] = rule
		ctx.HTML(200, "filterrule_new_or_edit")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func PostNewRule(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		ruleType := strings.TrimSpace(ctx.Req.Form.Get("ruletype"))
		ruleKey := strings.TrimSpace(ctx.Req.Form.Get("rulekey"))
		ruleValue := strings.TrimSpace(ctx.Req.Form.Get("rulevalue"))
		newRuleType, _ := strconv.Atoi(ruleType)
		rule := models.NewFilterRule(newRuleType, ruleKey, ruleValue)
		rule.Insert()
		ctx.Redirect("/admin/filterrules/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func EditRule(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		rule, _, _ := models.GetFilterRuleById(int64(Id))
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["rule"] = rule
		ctx.Data["type"] = "edit"
		ctx.Data["user"] = sess.Get("admin")
		ctx.HTML(200, "filterrule_new_or_edit")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func PostEditedRule(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		ruleType, _ := strconv.Atoi(strings.TrimSpace(ctx.Req.Form.Get("ruletype")))
		ruleKey := strings.TrimSpace(ctx.Req.Form.Get("rulekey"))
		ruleValue := strings.TrimSpace(ctx.Req.Form.Get("rulevalue"))
		models.EditFilterRuleById(int64(Id), ruleType, ruleKey, ruleValue)
		ctx.Redirect("/admin/filterrules/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DeleteRule(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id, _ := strconv.Atoi(ctx.Params(":id"))
		models.DeleteFilterRuleById(int64(id))
		ctx.Redirect("/admin/filterrules/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func ChangeRuleType(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id, _ := strconv.Atoi(ctx.Params(":id"))
		models.ConvertRuleType(int64(id))
		ctx.Redirect("/admin/filterrules/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

