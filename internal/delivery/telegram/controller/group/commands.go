package group

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func (h Handler) GetGroupCommand(c tele.Context) error {
	ctx := context.Background()
	chat, _, err := h.chats.Lookup(ctx, c.Chat().ID)
	if err != nil {
		h.logger.Error("cannot get chat", zap.Error(err))
		return c.Send("тех. ошибка")
	}

	if chat == nil {
		return c.Send("error: cannot get chat")
	}
	if chat.Group == nil {
		return c.Send("Группа для чата не установлена")
	}

	return c.Send(fmt.Sprintf("Для чата выбрана группа %s", *chat.Group))
}
