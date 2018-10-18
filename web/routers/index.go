package routers

import (
	"gopkg.in/macaron.v1"
)

func Index(ctx *macaron.Context) {
	ctx.Req.Header.Set("user", "anonymous")
	ctx.Redirect("/admin/index/")
}
