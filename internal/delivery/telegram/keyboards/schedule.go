package keyboards

import (
	"time"

	scheduleapi "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/telegram/callback"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/telegram/keyboard"
	tele "gopkg.in/telebot.v3"
)

var ScheduleCallback = callback.New("schedule", "day", "g")

func ScheduleButton(text string, day time.Time, group scheduleapi.Group) tele.Btn {
	return ScheduleCallback.MustTeleBtn(text, callback.M{
		"day": day.Format("2006-01-02"),
		"g":   group.String(),
	})
}

func ScheduleMarkup(today time.Time, group scheduleapi.Group) *tele.ReplyMarkup {
	b := keyboard.NewBuilder(2)
	b.Add(
		ScheduleButton("Пред. день", today.Add(-24*time.Hour), group),
		ScheduleButton("След. день", today.Add(24*time.Hour), group),
	)
	return b.Inline()
}
