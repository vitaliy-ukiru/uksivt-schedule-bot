package telegram

import (
	"context"
	"time"

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

func getStartDayOfWeek(tm time.Time) time.Time { //get monday 00:00:00
	weekday := time.Duration(tm.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	year, month, day := tm.Date()
	currentZeroDay := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return currentZeroDay.Add(-1 * (weekday - 1) * 24 * time.Hour)
}
