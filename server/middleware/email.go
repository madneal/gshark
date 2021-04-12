package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/madneal/gshark/global"
	"github.com/madneal/gshark/model"
	"github.com/madneal/gshark/model/request"
	"github.com/madneal/gshark/service"
	"github.com/madneal/gshark/utils"
	"go.uber.org/zap"
	"io/ioutil"
	"strconv"
	"time"
)

func ErrorToEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		var username string
		if claims, ok := c.Get("claims"); ok {
			waitUse := claims.(*request.CustomClaims)
			username = waitUse.Username
		} else {
			id, _ := strconv.Atoi(c.Request.Header.Get("x-user-id"))
			err, user := service.FindUserById(id)
			if err != nil {
				username = "Unknown"
			}
			username = user.Username
		}
		body, _ := ioutil.ReadAll(c.Request.Body)
		record := model.SysOperationRecord{
			Ip:     c.ClientIP(),
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
			Agent:  c.Request.UserAgent(),
			Body:   string(body),
		}
		now := time.Now()

		c.Next()

		latency := time.Now().Sub(now)
		status := c.Writer.Status()
		record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		str := "接收到的请求为" + record.Body + "\n" + "请求方式为" + record.Method + "\n" + "报错信息如下" + record.ErrorMessage + "\n" + "耗时" + latency.String() + "\n"
		if status != 200 {
			subject := username + "" + record.Ip + "调用了" + record.Path + "报错了"
			if err := utils.ErrorToEmail(subject, str); err != nil {
				global.GVA_LOG.Error("ErrorToEmail Failed, err:", zap.Any("err", err))
			}
		}
	}
}
