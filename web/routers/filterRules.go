package routers

import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/session"
	"strconv"
	"../../models"
	"strings"
	"github.com/go-macaron/csrf"
)

func ListFilterRules(ctx *macaron.Context, sess session.Store) {
	page := ctx.Params(":page")
	p, _ := strconv.Atoi(page)
	if p < 1 {
		p = 1
	}
	var pre int

	if p <= 1 {
		pre = 1
	} else {
		pre = p -1
	}

	next := p + 1

	if sess.Get("admin") != nil {
		rules, pages, _ := models.GetFilterRulesPage(p)
		pList := 0
		if (pages - p) > 10 {
			pList = p + 10
		} else {
			pList = pages
		}

		pageList := make([]int, 0)

		if pages <= 10 {
			for i:= 1; i < pList; i++ {
				pageList = append(pageList, i)
			}
		} else {
			if p <= 10 {
				for i := 1; i <= pList; i++ {
					pageList = append(pageList, i)
				}
			} else {
				t := p + 5
				if t > pages {
					t = pages
				}

				for i := p - 5; i <= t; i++ {
					pageList = append(pageList, i)
				}
			}
		}

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
		ctx.HTML(200, "filterrule_new")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func PostNewRule(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("amin") != nil {
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

