package routers

import (
	"github.com/madneal/gshark/models"
	"gopkg.in/macaron.v1"

	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"

	"strconv"
	"strings"
)

func ListTokens(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		tokens, _ := models.ListTokens()
		ctx.Data["tokens"] = tokens
		ctx.Data["role"] = sess.Get("user").(string)
		ctx.HTML(200, "tokens")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func NewTokens(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		ctx.Data["role"] = sess.Get("user").(string)
		ctx.HTML(200, "tokens_new")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoNewTokens(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		tokens := strings.TrimSpace(ctx.Req.Form.Get("tokens"))
		desc := strings.TrimSpace(ctx.Req.Form.Get("desc"))
		tokenType := strings.TrimSpace(ctx.Req.Form.Get("type"))
		githubToken := models.NewGithubToken(tokens, desc, tokenType)
		githubToken.Insert()
		ctx.Redirect("/admin/tokens/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func EditTokens(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		tokens, _, _ := models.GetTokenById(int64(Id))
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["tokens"] = tokens
		ctx.Data["admin"] = sess.Get("admin")
		ctx.Data["role"] = sess.Get("user").(string)
		ctx.HTML(200, "tokens_edit")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoEditTokens(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		tokens := strings.TrimSpace(ctx.Req.Form.Get("tokens"))
		desc := strings.TrimSpace(ctx.Req.Form.Get("desc"))
		tokenType := strings.TrimSpace(ctx.Req.Form.Get("type"))
		models.EditTokenById(int64(Id), tokens, desc, tokenType)
		ctx.Redirect("/admin/tokens/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DeleteTokens(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		models.DeleteTokenById(int64(Id))
		ctx.Redirect("/admin/tokens/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}
