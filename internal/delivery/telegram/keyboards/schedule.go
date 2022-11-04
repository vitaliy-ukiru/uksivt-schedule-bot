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
		ScheduleButton("Пред. день", addDays(today, -1, time.Sunday), group),
		ScheduleButton("След. день", addDays(today, +1, time.Sunday), group),
	)

	return b.Inline()
}

func addDays(t time.Time, days int, skip time.Weekday) time.Time {
	duration := time.Duration(days) * 24 * time.Hour
	t = t.Add(duration)
	if t.Weekday() == skip {
		t = t.Add(duration)
	}

	return t
}
