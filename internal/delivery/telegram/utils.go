package telegram

import (
	"context"

	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat"
	"go.uber.org/zap"
)

type Chat struct {
	*chat.Chat
	Status chat.LookupStatus
}

func (h Handler) getChat(tgID int64) *Chat {
	chatObj, status, err := h.uc.Lookup(context.TODO(), tgID)
	if err != nil {
		h.logger.Error("cannot get chat", zap.Error(err))
		return nil
	}
	return &Chat{Chat: chatObj, Status: status}
}
