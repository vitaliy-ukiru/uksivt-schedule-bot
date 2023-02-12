package cron

import (
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/telegram/keyboard"
	tele "gopkg.in/telebot.v3"
)

//const FlagsCallback = "cron_c_flag"

var BackBtn = keyboard.CallbackButton("Назад", "cron_c_back")

var (
	CancelBtn = keyboard.CallbackButton("Отмена", "cron_cancel")
	AcceptBtn = keyboard.CallbackButton("Создать", "cron_c_accept")
)

func AcceptMarkup() *tele.ReplyMarkup {
	return keyboard.
		NewBuilder(1).
		Add(CancelBtn, AcceptBtn).
		Inline()
}