package sauth

import (
	"github.com/casbin/casbin"
	"gopkg.in/macaron.v1"
	"net/http"
)


func Authorizer(e *casbin.Enforcer) macaron.Handler {
	return func(res http.ResponseWriter, req *http.Request, c *macaron.Context) {
		user := c.GetCookie("user")
		if user == "" {
			user = "anonymous"
		}
		method := req.Method
		path := req.URL.Path
		if !e.Enforce(user, path, method) {
			accessDenied(res)
			return
		}
	}
}

func accessDenied(res http.ResponseWriter) {
	http.Error(res, "Access Denied", http.StatusForbidden)
}
