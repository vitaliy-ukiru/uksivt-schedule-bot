package cron

import (
	"strconv"

	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/scheduler"
)

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
