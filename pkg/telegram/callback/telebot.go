package callback

import (
	"strings"

	tele "gopkg.in/telebot.v3"
)

func (d Data) ParseTele(callback *tele.Callback) (M, error) {
	prefix, parts := callback.Unique, strings.Split(callback.Data, TelebotSeparator)
	return d.union(prefix, parts)
}

func (d Data) TeleBtn(text string, kwParts M, parts ...string) (btn tele.Btn, err error) {
	data, err := d.buildData(kwParts, parts)
	if err != nil {
		return
	}

	btn.Unique = d.Prefix
	btn.Text = text
	btn.Data = strings.Join(data, TelebotSeparator)

	return
}

func (d Data) MustTeleBtn(text string, kwParts M, parts ...string) tele.Btn {
	btn, err := d.TeleBtn(text, kwParts, parts...)
	if err != nil {
		panic(err)
	}
	return btn
}

func (d Data) CallbackUnique() string {
	return "\f" + d.Prefix
}
