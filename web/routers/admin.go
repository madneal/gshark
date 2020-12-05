package routers

import (
	"github.com/madneal/gshark/logger"
	"github.com/madneal/gshark/models"

	"gopkg.in/macaron.v1"

	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"

	"strconv"
	"strings"
)

func AdminIndex(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		ctx.Redirect("/admin/reports/github/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func Login(ctx *macaron.Context, x csrf.CSRF) {
	ctx.Data["csrf_token"] = x.GetToken()
	ctx.HTML(200, "login")
}

func DoLogin(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	username := ctx.Req.Form.Get("username")
	password := ctx.Req.Form.Get("password")
	has, role, err := models.Auth(username, password)
	if err == nil && has {
		sess.Set("user", role)
		sess.Set("admin", username)
		ctx.SetCookie("user", role)
		ctx.Redirect("/admin/index/")
	} else {
		ctx.Data["login_error"] = "用户名或者密码错误！"
		ctx.HTML(200, "login")
	}
}

func DoLogout(ctx *macaron.Context, sess session.Store) {
	sess.GC()
	ctx.SetCookie("user", "anonymous")
	ctx.Redirect("/admin/login")
}

func ListUsers(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		users, _ := models.ListAdmins()
		ctx.Data["users"] = users
		ctx.Data["role"] = sess.Get("user").(string)
		ctx.HTML(200, "users")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func NewUser(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		ctx.Data["role"] = sess.Get("user").(string)
		ctx.HTML(200, "users_new")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoNewUser(ctx *macaron.Context, sess session.Store) {
	ctx.Req.ParseForm()
	if sess.Get("admin") != nil {
		username := strings.TrimSpace(ctx.Req.Form.Get("username"))
		password := strings.TrimSpace(ctx.Req.Form.Get("password"))
		role := ctx.Req.Form.Get("role")
		admin := models.NewAdmin(username, password, role)
		admin.Insert()
		ctx.Redirect("/admin/users/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func EditUser(ctx *macaron.Context, sess session.Store, x csrf.CSRF) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		admin, _, _ := models.GetAdminById(int64(Id))
		ctx.Data["csrf_token"] = x.GetToken()
		ctx.Data["user"] = admin
		ctx.Data["admin"] = sess.Get("admin")
		ctx.Data["role"] = sess.Get("user").(string)
		ctx.HTML(200, "users_edit")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DoEditUser(ctx *macaron.Context, sess session.Store) {
	err := ctx.Req.ParseForm()
	if err != nil {
		logger.Log.Error(err)
	}
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		username := strings.TrimSpace(ctx.Req.Form.Get("username"))
		password := strings.TrimSpace(ctx.Req.Form.Get("password"))
		role := ctx.Req.Form.Get("role")
		err := models.EditAdminById(int64(Id), username, password, role)
		if err != nil {
			logger.Log.Error(err)
		}
		ctx.Redirect("/admin/users/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DeleteUser(ctx *macaron.Context, sess session.Store) {
	if sess.Get("admin") != nil {
		id := ctx.Params(":id")
		Id, _ := strconv.Atoi(id)
		err := models.DeleteAdminById(int64(Id))
		if err != nil {
			logger.Log.Error(err)
		}
		ctx.Redirect("/admin/users/list/")
	} else {
		ctx.Redirect("/admin/login/")
	}
}
