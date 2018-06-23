package main

import (
	"bytes"
	"github.com/tucnak/telebot"
)

type Integration interface {
	Platform() string
	Send(sqlEvent) error
}

type Telegram struct {
	Token  string
	ChatID string
}

func (t *Telegram) Platform() string {
	return "telegram"
}

func (t *Telegram) Send(event sqlEvent) error {
	bot, err := telebot.NewBot(telebot.Settings{
		Token: t.Token,
	})
	if err != nil {
		return err
	}
	chat, err := bot.ChatByID(t.ChatID)
	if err != nil {
		return err
	}
	_, err = bot.Send(chat, t.TXT(event), telebot.Silent)
	return err
}

func (t *Telegram) TXT(event sqlEvent) string {
	buf := &bytes.Buffer{}
	buf.WriteString("[")
	buf.WriteString(event.Network)
	buf.WriteString("]")
	buf.WriteString(" ")
	buf.WriteString(event.ContractName)
	buf.WriteString(" ")
	buf.WriteString(event.EventName)
	buf.WriteString("(")
	for i := range event.Arguments {
		if i != 0 {
			buf.WriteString(",")
		}
		buf.WriteString(event.Arguments[i])
	}
	buf.WriteString(")")
	return buf.String()
}
