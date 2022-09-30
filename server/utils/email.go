package utils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/madneal/gshark/global"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"strings"
)

func Email(subject string, body string) error {
	to := strings.Split(global.GVA_CONFIG.Email.To, ",")
	return send(to, subject, body)
}

func ErrorToEmail(subject string, body string) error {
	to := strings.Split(global.GVA_CONFIG.Email.To, ",")
	if to[len(to)-1] == "" { // 判断切片的最后一个元素是否为空,为空则移除
		to = to[:len(to)-1]
	}
	return send(to, subject, body)
}

func EmailSend(subject string, body string) error {
	to := []string{global.GVA_CONFIG.Email.From}
	return send(to, subject, body)
}

func BotSend(content string) error {
	url := global.GVA_CONFIG.Wechat.Url
	if url == "" {
		err := errors.New("url is empty")
		return err
	}
	jsonStr := []byte(fmt.Sprintf(`{"msgtype": "markdown", "markdown":{"content":"%s"}}`, content))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	c := &http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	return err
}

//@author: [maplepie](https://github.com/maplepie)
//@function: send
//@description: Email发送方法
//@param: subject string, body string
//@return: error

func send(to []string, subject string, body string) error {
	from := global.GVA_CONFIG.Email.From
	nickname := global.GVA_CONFIG.Email.Nickname
	secret := global.GVA_CONFIG.Email.Secret
	host := global.GVA_CONFIG.Email.Host
	port := global.GVA_CONFIG.Email.Port
	isSSL := global.GVA_CONFIG.Email.IsSSL

	auth := smtp.PlainAuth("", from, secret, host)
	e := email.NewEmail()
	if nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", nickname, from)
	} else {
		e.From = from
	}
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)
	var err error
	hostAddr := fmt.Sprintf("%s:%d", host, port)
	if isSSL {
		err = e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: host})
	} else {
		err = e.Send(hostAddr, auth)
	}
	return err
}
