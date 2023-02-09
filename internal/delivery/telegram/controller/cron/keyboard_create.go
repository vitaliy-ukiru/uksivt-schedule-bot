package cron

import (
	"time"

	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/scheduler"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/telegram/callback"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/telegram/keyboard"
	tele "gopkg.in/telebot.v3"
)

var SelectTimeCallback = callback.New("cron_c", "time", "user")

var (
	PMButton = keyboard.CallbackButton("После полудня", "cron_c_pm")
	AMButton = keyboard.CallbackButton("До полудня", "cron_c_am")
)

func AMTimesMarkup(userId string, period time.Duration) *tele.ReplyMarkup {
	b := keyboard.NewBuilder(4)
	var t time.Time
	for t.Hour() < 12 {
		tStr := t.Format("15:04")
		b.Insert(SelectTimeCallback.MustTeleBtn(
			tStr,
			callback.M{
				"time": tStr,
				"user": userId,
			},
		))
		t = t.Add(period)
	}
	b.OneButtonRow(PMButton)
	return b.Inline()
}

func PMTimesMarkup(userId string, period time.Duration) *tele.ReplyMarkup {
	b := keyboard.NewBuilder(4)

	var t time.Time
	t = t.Add(12 * time.Hour)

	for t.Hour() != 0 {
		tStr := t.Format("15:04")
		b.Insert(SelectTimeCallback.MustTeleBtn(
			tStr,
			callback.M{
				"time": tStr,
				"user": userId,
			},
		))
		t = t.Add(period)

	}
	b.OneButtonRow(AMButton)
	return b.Inline()
}

//const FlagsCallback = "cron_c_flag"

var FlagsCallback = callback.New("cron_c_flag", "mode")

var AcceptFlags = keyboard.CallbackButton("Подтвердить", "cron_c_flag_acc")

func FlagsMarkup(set scheduler.FlagSet) *tele.ReplyMarkup {
	b := keyboard.NewBuilder(1)

	for i, mode := range FlagModes {
		b.OneButtonRow(FlagsCallback.MustTeleBtn(
			mode.FormatText(set, i, "✅"),
			nil,
			mode.Callback,
		))
	}
	b.OneButtonRow(AcceptFlags)

	return b.Inline()

}

var BackBtn = keyboard.CallbackButton("Назад", "cron_c_back")

var (
	CancelBtn = keyboard.CallbackButton("Отмена", "cron_cancel")
	AcceptBtn = keyboard.CallbackButton("Создать", "cron_c_accept")
)

func AcceptMarkup() *tele.ReplyMarkup {
	return keyboard.
		NewBuilder(1).
		Add(CancelBtn,
			AcceptBtn).
		Inline()
}
