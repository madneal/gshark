package service

import (
	"github.com/madneal/gshark/utils"
)

func EmailTest() (err error) {
	subject := "[GShark] email test"
	body := "This is a test email from GShark to verify the email send function."
	err = utils.EmailSend(subject, body)
	return err
}

func BotTest() (err error) {
	content := "test"
	err = utils.BotSend(content)
	return err
}
