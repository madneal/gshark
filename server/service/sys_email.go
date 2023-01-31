package service

import (
	"github.com/madneal/gshark/utils"
)

func EmailTest() (err error) {
	subject := "test"
	body := "test"
	err = utils.EmailSend(subject, body)
	return err
}

func BotTest() (err error) {
	content := "test"
	err = utils.BotSend(content)
	return err
}
