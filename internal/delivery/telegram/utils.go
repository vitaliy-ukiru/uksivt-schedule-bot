package telegram

import (
	"time"

	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat"
	tele "gopkg.in/telebot.v3"
)

func getChat(c tele.Context) *chat.Chat {
	obj := c.Get(ChatKey)
	if obj == nil {
		return nil
	}

	chat2, ok := obj.(*chat.Chat)
	if !ok {
		return nil
	}
	return chat2
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
