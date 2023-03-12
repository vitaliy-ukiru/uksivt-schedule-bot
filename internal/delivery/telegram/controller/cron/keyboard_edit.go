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
	SelectEditTitle  = keyboard.Text("Название")
	SelectEditTime   = keyboard.Text("Время")
	SelectEditFlags  = keyboard.Text("Опции")
	DoneEditing      = keyboard.Text("Сохранить")
	CancelEditingBtn = keyboard.Text("Закрыть")
)

func SelectEditingFieldMarkup() *tele.ReplyMarkup {
	m := keyboard.NewBuilderBuffer(3, 2).
		Add(SelectEditTitle, SelectEditTime, SelectEditFlags).
		Add(DoneEditing, CancelEditingBtn).
		Reply()
	m.ResizeKeyboard = true
	return m
}
