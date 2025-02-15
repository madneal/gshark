package utils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
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

func send(to []string, subject string, body string) error {
	from := global.GVA_CONFIG.Email.From
	secret := global.GVA_CONFIG.Email.Secret
	host := global.GVA_CONFIG.Email.Host
	port := global.GVA_CONFIG.Email.Port
	smtpServer := fmt.Sprintf("%s:%d", host, port)

	auth := smtp.PlainAuth("", from, secret, host)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	conn, err := tls.Dial("tcp", smtpServer, tlsConfig)
	if err != nil {
		return fmt.Errorf("dial err: %v", err)
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp client err: %v", err)
	}
	defer c.Quit()

	if err = c.Auth(auth); err != nil {
		return fmt.Errorf("auth err: %v", err)
	}

	if err = c.Mail(from); err != nil {
		return fmt.Errorf("mail err: %v", err)
	}
	if err = c.Rcpt(to[0]); err != nil {
		return fmt.Errorf("rcpt err: %v", err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("data err: %v", err)
	}
	defer w.Close()

	msg := []byte("From: Sender Name <" + from + ">\r\n" +
		"To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")
	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("write body err: %v", err)
	}
	return nil
}
