package cron

import (
	"strconv"
	"time"

	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/scheduler"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/telegram/callback"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/telegram/keyboard"
	tele "gopkg.in/telebot.v3"
)

var SelectTimeCallback = callback.New("cron_c", "time", "user")

var (
	PMButton = keyboard.CallbackButton("После полудня", "cron_c_pm")
	AMButton = keyboard.CallbackButton("До полудня", "cron_c_am")
)

func TimesMarkupAM(userId string, period time.Duration) *tele.ReplyMarkup {
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

func TimesMarkupPM(userId string, period time.Duration) *tele.ReplyMarkup {
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

var FlagModes = []FlagMode{
	{
		Mode:     scheduler.NextDay,
		Callback: "NextDay",
		Text:     "На след. день",
	},
	{
		Mode:     scheduler.Full,
		Callback: "Full",
		Text:     "Полное расписание",
	},
	{
		Mode:     scheduler.OnlyIfHaveReplaces,
		Callback: "OnlyIfReplaces",
		Text:     "Только если есть замены",
	},
	{
		Mode:     scheduler.ReplacesAlways,
		Callback: "ReplacesAlways",
		Text:     "Только замены (иначе увед.)",
	},
	{
		Mode:     scheduler.FullOnlyIfReplaces,
		Callback: "FullOnReplaces",
		Text:     "Полное, только если замены",
	},
}

type FlagMode struct {
	Mode     scheduler.FlagSet
	Callback string
	Text     string
}

func FlagSetFromCallback(s string) (scheduler.FlagSet, bool) {
	for _, mode := range FlagModes {
		if mode.Callback == s {
			return mode.Mode, true
		}
	}
	return 0, false
}

func (c FlagMode) FormatText(f scheduler.FlagSet, i int, add string) string {
	text := c.Text
	if i > -1 {
		text = strconv.Itoa(i) + ". " + text
	}
	if f.Has(c.Mode) {
		text += " " + add
	}
	return text
}

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
