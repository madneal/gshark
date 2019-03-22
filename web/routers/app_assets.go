package routers

import (
	"gopkg.in/macaron.v1"
	"github.com/go-macaron/session"
	"strconv"
	"github.com/neal1991/gshark/util/common"
	"github.com/neal1991/gshark/models"
	"github.com/neal1991/gshark/vars"
	"os"
	"log"
	"io"
	sha2562 "crypto/sha256"
	"encoding/hex"
	"fmt"
)

type HashForm struct {
	MD5     string
}

func ListAppAssets(ctx *macaron.Context, sess session.Store) {
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
		ctx.HTML(200, "app_assets")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DetectApp(ctx *macaron.Context, sess session.Store)  {
	if sess.Get("admin") != nil {
		hash := ctx.Req.Form.Get("hash")
		fmt.Println(hash)
		ctx.HTML(200, "app_detect")
	} else {
		ctx.Redirect("/admin/login/")
	}
}

func DetectAppExists(filepath string) bool {
	hash := GenerateFileHash(filepath)
	return models.Detect(hash)
}

//generate sha256 of file
func GenerateFileHash(filepath string) (hashResult string) {
	hash := sha2562.New()
	f, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if _, err := io.Copy(hash, f); err != nil {
		log.Fatal(err)
	}
	hashResult = hex.EncodeToString(hash.Sum(nil))
	return hashResult
}

