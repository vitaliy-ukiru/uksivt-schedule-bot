package keyboard

import tele "gopkg.in/telebot.v3"

var markup = &tele.ReplyMarkup{}

//goland:noinspection GoUnusedGlobalVariable
var (
	Text     = markup.Text
	Contact  = markup.Contact
	Location = markup.Location
	Poll     = markup.Poll

	CallbackButton  = markup.Data
	URLButton       = markup.URL
	QueryButton     = markup.Query
	QueryChatButton = markup.QueryChat

	Row = markup.Row
)
