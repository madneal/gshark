package web

import (
	"fmt"
	"github.com/casbin/casbin"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	"github.com/neal1991/gshark/logger"
	"github.com/neal1991/gshark/sauth"
	"github.com/neal1991/gshark/vars"
	"github.com/neal1991/gshark/web/routers"
	"github.com/urfave/cli"
	"gopkg.in/macaron.v1"
	"html/template"
	"net/http"
	"runtime"
	"strings"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func RunWeb(ctx *cli.Context) {

	if ctx.IsSet("debug") {
		vars.DEBUG_MODE = ctx.Bool("debug")
	}
	if ctx.IsSet("host") {
		vars.HTTP_HOST = ctx.String("host")
	}
	if ctx.IsSet("port") {
		vars.HTTP_PORT = ctx.Int("port")
	}

	m := macaron.Classic()
	e := casbin.NewEnforcer("./conf/auth_model.conf", "./conf/policy.csv")
	m.Use(sauth.Authorizer(e))

	m.Use(macaron.Renderer(
		macaron.RenderOptions{
			Directory:  "templates",
			Extensions: []string{".tmpl", ".html"},
			Funcs: []template.FuncMap{map[string]interface{}{
				"Replace": func(str *string) string {
					t := strings.Replace(*str, "\n", "<br>", -1)
					return t
				},
				"Split": func(str *string) []string {
					return strings.Split(*str, ",")
				},
				"unescaped": func(x string) interface{} { return template.HTML(x) },
			}},
			Delims:          macaron.Delims{"{{", "}}"},
			Charset:         "UTF-8",
			IndentJSON:      true,
			IndentXML:       true,
			//PrefixJSON:      []byte("macaron"),
			PrefixXML:       []byte("macaron"),
			HTMLContentType: "text/html",
		}))

	m.Use(session.Sessioner(session.Options{
		// GC 执行时间间隔，默认为 3600 秒
		Gclifetime: 3600,
		// 最大生存时间，默认和 GC 执行时间间隔相同
		Maxlifetime: 3600,
	}))
	m.Use(csrf.Csrfer())
	m.Use(cache.Cacher())
	m.Get("/", routers.Index)
	m.Group("/admin", func() {
		m.Get("", routers.AdminIndex)
		m.Get("/index/", routers.AdminIndex)
		m.Get("/login/", routers.Login)
		m.Post("/login/", routers.DoLogin)
		m.Get("/logout", routers.DoLogout)

		m.Group("/users/", func() {
			m.Get("", routers.ListUsers)
			m.Get("/list/", routers.ListUsers)
			m.Get("/new/", routers.NewUser)
			m.Post("/new/", routers.DoNewUser)
			m.Get("/edit/:id", routers.EditUser)
			m.Post("/edit/:id", routers.DoEditUser)
			m.Get("/del/:id", routers.DeleteUser)
		})

		m.Group("/assets/", func() {
			m.Get("", routers.ListAssets)
			m.Get("/list/", routers.ListAssets)
			m.Get("/list/:page", routers.ListAssets)
			m.Get("/new/", routers.NewAssets)
			m.Post("/new/", routers.DoNewAssets)
			m.Get("/edit/:id", routers.EditAssets)
			m.Post("/edit/:id", routers.DoEditAssets)
			m.Get("/del/:id", routers.DeleteAssets)
			m.Get("/del_all/", routers.DeleteAllAssets)
		})

		m.Group("/app/", func() {
			m.Get("", routers.ListAppAssets)
			m.Get("/list", routers.ListAppAssets)
			m.Get("/detect/", routers.DetectApp)
			m.Get("/appid/", routers.GetAppAsset)
			m.Get("/new", routers.NewAppAsset)
			m.Post("/new", routers.DoNewAppAsset)
		})

		m.Group("/tokens/", func() {
			m.Get("", routers.ListTokens)
			m.Get("/list/", routers.ListTokens)
			m.Get("/new/", routers.NewTokens)
			m.Post("/new/", routers.DoNewTokens)
			m.Get("/edit/:id", routers.EditTokens)
			m.Post("/edit/:id", routers.DoEditTokens)
			m.Get("/del/:id", routers.DeleteTokens)
		})

		m.Group("/rules/", func() {
			m.Get("", routers.ListRules)
			m.Get("/list/", routers.ListRules)
			m.Get("/list/:page", routers.ListRules)
			m.Get("/new/", routers.NewRules)
			m.Post("/new/", routers.DoNewRules)
			m.Get("/edit/:id", routers.EditRules)
			m.Post("/edit/:id", routers.DoEditRules)
			m.Get("/del/:id", routers.DeleteRules)
			m.Get("/enable/:id", routers.EnableRules)
			m.Get("/disable/:id", routers.DisableRules)
		})

		m.Group("/filterrules/", func() {
			m.Get("", routers.ListFilterRules)
			m.Get("/list/", routers.ListFilterRules)
			m.Get("/list/:page", routers.ListFilterRules)
			m.Get("/new/", routers.NewRule)
			m.Post("/new/", routers.PostNewRule)
			m.Get("/edit/:id", routers.EditRule)
			m.Post("/edit/:id", routers.PostEditedRule)
			m.Get("/del/:id", routers.DeleteRule)
		})

		m.Group("/repos/", func() {
			m.Get("", routers.ListRepos)
			m.Get("/list/", routers.ListRepos)
			m.Get("/list/:page", routers.ListRepos)
			m.Get("/enable/:id", routers.EnableRepo)
			m.Get("/disable/:id", routers.DisableRepo)
			m.Get("/del/", routers.DeleteAllRepo)
		})

		m.Group("/reports/", func() {
			m.Get("/github/", routers.ListGithubSearchResult)
			m.Get("/github/:page", routers.ListGithubSearchResult)
			m.Get("/github/confirm/:id", routers.ConfirmReportById)
			m.Get("/github/cancel/:id", routers.CancelReportById)
			m.Get("/github/disable_repo/:id", routers.DisableRepoById)
			m.Get("/github/report_detail/:id", routers.GetDetailedReportById)
			m.Get("/github/query/:status", routers.ListGithubSearchResultByStatus)
			m.Get("/app/", routers.ListAppSearchResult)
			m.Get("/app/:page", routers.ListAppSearchResult)
			m.Get("/app/confirm/:id", routers.ConfirmReportById)
			m.Get("/app/cancel/:id", routers.CancelReportById)
			m.Get("/app/query/:status", routers.ListAppSearchResultByStatus)
		})
	})

	logger.Log.Printf("Server is running on %s", fmt.Sprintf("%v:%v", vars.HTTP_HOST, vars.HTTP_PORT))
	logger.Log.Println(http.ListenAndServe(fmt.Sprintf("%v:%v", vars.HTTP_HOST, vars.HTTP_PORT), m))
}
