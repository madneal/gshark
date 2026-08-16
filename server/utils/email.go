package utils

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/madneal/gshark/global"
)

func EmailSend(subject string, body string) error {
	to := splitRecipients(global.GVA_CONFIG.Email.To)
	if len(to) == 0 {
		to = []string{global.GVA_CONFIG.Email.From}
	}
	return send(to, subject, body)
}

func splitRecipients(value string) []string {
	raw := strings.FieldsFunc(value, func(r rune) bool {
		return r == ',' || r == ';' || r == '\n'
	})
	recipients := make([]string, 0, len(raw))
	for _, recipient := range raw {
		if recipient = strings.TrimSpace(recipient); recipient != "" {
			recipients = append(recipients, recipient)
		}
	}
	return recipients
}

func BotSend(content string) error {
	url := global.GVA_CONFIG.Wechat.Url
	if url == "" {
		return errors.New("url is empty")
	}

	payload := map[string]any{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": content,
		},
	}
	jsonStr, err := json.Marshal(payload)
	if err != nil {
		return err
	}

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

	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))
	return nil
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
	for _, recipient := range to {
		if err = c.Rcpt(recipient); err != nil {
			return fmt.Errorf("rcpt err: %v", err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("data err: %v", err)
	}
	defer w.Close()

	msg := []byte("From: Sender Name <" + from + ">\r\n" +
		"To: " + strings.Join(to, ", ") + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body + "\r\n")
	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("write body err: %v", err)
	}
	return nil
}
