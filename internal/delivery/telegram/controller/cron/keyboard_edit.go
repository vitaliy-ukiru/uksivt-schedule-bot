package cron

import (
	"strconv"
	"strings"

	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/scheduler"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/telegram/callback"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/telegram/keyboard"
	tele "gopkg.in/telebot.v3"
)

var SelectCronEditCallback = callback.New("cron_e_select", "id")

func SelectCronsMarkup(crons []scheduler.CronJob) *tele.ReplyMarkup {
	b := keyboard.NewBuilderBuffer(1, len(crons))
	for _, cron := range crons {
		b.OneButtonRow(SelectCronEditCallback.MustTeleBtn(
			strings.Join(
				[]string{
					cron.Title,
					cron.At.Format("15:04"),
				},
				" | ",
			),
			nil,
			strconv.FormatInt(cron.ID, 10),
		))
	}

	return b.Inline()
}

var (
	SelectEditTitle  = keyboard.CallbackButton("Название", "cron_e_title")
	SelectEditTime   = keyboard.CallbackButton("Время", "cron_e_time")
	SelectEditFlags  = keyboard.CallbackButton("Опции", "cron_e_flags")
	DoneEditing      = keyboard.CallbackButton("Сохранить", "cron_e_done")
	CancelEditingBtn = keyboard.CallbackButton("Закрыть", "cron_e_cancel")
)

var SelectEditingFieldMarkup = keyboard.
	NewBuilderBuffer(3, 2).
	Add(SelectEditTitle, SelectEditTime, SelectEditFlags).
	Add(DoneEditing, CancelEditingBtn).
	Inline()
