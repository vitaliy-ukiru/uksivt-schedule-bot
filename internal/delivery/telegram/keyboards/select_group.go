package keyboards

import (
	"strconv"

	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/telegram/keyboard"
	tele "gopkg.in/telebot.v3"
)

const SGCallback = "select_g"

var (
	CancelBtn = keyboard.CallbackButton("Отмена", "select_g_cancel")
	AcceptBtn = keyboard.CallbackButton("Подтвердить", "select_g_accept")
	//BackBtn   = keyboard.CallbackButton("Назад", "select.g.back")
)

func SelectYear(years []int) *tele.ReplyMarkup {
	b := keyboard.NewBuilder(1)
	for _, year := range years {
		yearStr := strconv.Itoa(year)
		btn := keyboard.CallbackButton(yearStr+" год", SGCallback, yearStr)
		b.OneButtonRow(btn)
	}
	b.OneButtonRow(CancelBtn)
	return b.Inline()
}

func SelectSpec(year int, specs []string) *tele.ReplyMarkup {
	yearStr := strconv.Itoa(year)
	const rowSize = 3.0
	b := keyboard.NewBuilder(rowSize)
	for _, spec := range specs {
		b.Insert(keyboard.CallbackButton(yearStr+spec, SGCallback, spec))
	}
	b.OneButtonRow(CancelBtn)

	return b.Inline()
}

func SelectNumber(year int, spec string, numbers []int) *tele.ReplyMarkup {
	groupStr := strconv.Itoa(year) + spec + "-"
	b := keyboard.NewBuilder(1)
	for _, number := range numbers {
		numberStr := strconv.Itoa(number)
		btn := keyboard.CallbackButton(groupStr+numberStr, SGCallback, numberStr)
		b.OneButtonRow(btn)
	}
	b.OneButtonRow(CancelBtn)
	return b.Inline()
}

func AcceptMarkup() *tele.ReplyMarkup {
	b := keyboard.NewBuilder(1)
	b.Add(
		AcceptBtn,
		CancelBtn,
	)
	return b.Inline()

}
