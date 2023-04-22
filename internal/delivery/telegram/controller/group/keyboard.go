package group

import (
	"strconv"

	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/telegram/keyboard"
	tele "gopkg.in/telebot.v3"
)

const SGCallback = "select_g"

var (
	AcceptBtn = keyboard.CallbackButton("Подтвердить", SGCallback+"_accept")
	CancelBtn = keyboard.CallbackButton("Отмена", SGCallback+"_cancel")
	BackBtn   = keyboard.CallbackButton("Назад", SGCallback+"_back")

	CancelRow = tele.Row{BackBtn, CancelBtn}
)

func SelectYearMarkup(years []int) *tele.ReplyMarkup {
	b := keyboard.NewBuilderBuffer(1, len(years)+1)
	for _, year := range years {
		yearStr := strconv.Itoa(year)
		btn := keyboard.CallbackButton(yearStr+" год", SGCallback, yearStr)
		b.OneButtonRow(btn)
	}
	b.OneButtonRow(CancelBtn)
	return b.Inline()
}

func SelectSpecMarkup(year int, specs []string) *tele.ReplyMarkup {
	yearStr := strconv.Itoa(year)
	const rowSize = 3.0
	b := keyboard.NewBuilderBuffer(rowSize, len(specs)/rowSize+1)
	for _, spec := range specs {
		b.Insert(keyboard.CallbackButton(yearStr+spec, SGCallback, spec))
	}
	b.Row(CancelRow)

	return b.Inline()
}

func SelectNumberMarkup(year int, spec string, numbers []int) *tele.ReplyMarkup {
	groupStr := strconv.Itoa(year) + spec + "-"
	b := keyboard.NewBuilderBuffer(1, len(numbers)+1)
	for _, number := range numbers {
		numberStr := strconv.Itoa(number)
		btn := keyboard.CallbackButton(groupStr+numberStr, SGCallback, numberStr)
		b.OneButtonRow(btn)
	}
	b.Row(CancelRow)
	return b.Inline()
}

func AcceptMarkup() *tele.ReplyMarkup {
	b := keyboard.NewBuilderBuffer(1, 2)
	b.OneButtonRow(AcceptBtn)
	b.Row(CancelRow)
	return b.Inline()

}
