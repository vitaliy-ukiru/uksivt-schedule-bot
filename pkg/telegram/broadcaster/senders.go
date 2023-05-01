package broadcaster

import (
	"strings"
	"text/template"

	tele "gopkg.in/telebot.v3"
)

// TextSender sends text as string.
type TextSender string

func (t TextSender) Send(bot *tele.Bot, to tele.Recipient, opts *tele.SendOptions) (*tele.Message, error) {
	return bot.Send(to, string(t), opts)
}

// MessageSender copies message.
type MessageSender struct {
	*tele.Message
}

func NewMessageSender(message *tele.Message) *MessageSender {
	return &MessageSender{Message: message}
}

func (m *MessageSender) Send(bot *tele.Bot, to tele.Recipient, opts *tele.SendOptions) (*tele.Message, error) {
	return bot.Copy(to, m.Message, opts)
}

// TemplateSender executes template for each recipient with parameters from Getter
type TemplateSender struct {
	Template *template.Template
	Getter   func(recipient tele.Recipient) any // data getter for current recipient
}

func NewTemplateSender(tmpl *template.Template, getter func(recipient tele.Recipient) any) *TemplateSender {
	return &TemplateSender{Template: tmpl, Getter: getter}
}

func (t TemplateSender) Send(bot *tele.Bot, to tele.Recipient, opts *tele.SendOptions) (*tele.Message, error) {
	var sb strings.Builder
	if err := t.Template.Execute(&sb, t.Getter(to)); err != nil {
		return nil, err
	}

	return bot.Send(to, sb.String(), opts)
}
