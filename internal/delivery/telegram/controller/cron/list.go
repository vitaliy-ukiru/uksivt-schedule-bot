package cron

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) ListCrons(c tele.Context) error {
	chat, _, err := h.chats.Lookup(context.TODO(), c.Chat().ID)
	if err != nil {
		return err
	}

	crons, err := h.crons.ForChat(context.TODO(), chat.ID)
	if err != nil {
		return c.Send("error: " + err.Error() + fmt.Sprintf("\n\n%+v\n", err))
	}

	for _, cron := range crons {
		err := c.Send(fmt.Sprintf(
			"Title: %s\nAt:%s\nFlags:%s\n",
			cron.Title,
			cron.At.Format("15:04"),
			flagString(cron.Flags, "; "),
		))
		if err != nil {
			h.logger.Error("cannot send message", zap.Error(err))
			return c.Send("Произошла ошибка")
		}
	}
	return nil

}
