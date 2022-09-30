package telegram

import (
	"context"
	"fmt"
	"strconv"

	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram/keyboards"
	scheduleapi "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
	tele "gopkg.in/telebot.v3"
)
import fsm "github.com/vitaliy-ukiru/fsm-telebot"

var (
	//SelectYear   = groupSG.New("year")
	groupSG = fsm.NewStateGroup("s.g")

	SelectSpec   = groupSG.New("spec")
	SelectNumber = groupSG.New("num")
	AcceptGroup  = groupSG.New("accept")
)

func (h Handler) SGCommand(c tele.Context, _ fsm.FSMContext) error {
	//chat, err := h.uc.ByTelegramID(context.TODO(), c.Chat().ID)
	//if err != nil {
	//	return c.Send("cannot get chat: " + err.Error())
	//}
	//state.Set(SelectYear)
	markup := keyboards.SelectYear(h.groups.Years())
	return c.Send("Выберите год поступления:", markup)
}

func (h Handler) SGYearCallback(c tele.Context, state fsm.FSMContext) error {
	yearStr := c.Data()
	year, _ := strconv.Atoi(yearStr)
	//if err != nil {
	//	return err
	//}
	state.Set(SelectSpec)
	_ = state.Update("g", scheduleapi.Group{Year: year})
	specs := h.groups.Specs(year)
	markup := keyboards.SelectSpec(year, specs)

	return c.EditOrSend("Выберите специальность", markup)
}

func (h Handler) SGSpecCallback(c tele.Context, state fsm.FSMContext) error {
	//year, ok := state.MustGet("year").(int)
	g, ok := state.MustGet("g").(scheduleapi.Group)
	if !ok {
		_ = state.Finish(true)
		return c.Send("invalid data, aborting")
	}
	g.Spec = c.Data()

	_ = state.Update("g", g)
	state.Set(SelectNumber)
	numbers := h.groups.Numbers(g.Year, g.Spec)
	markup := keyboards.SelectNumber(g.Year, g.Spec, numbers)
	return c.EditOrSend("Выберите группу", markup)
}

func (h Handler) SGNumCallback(c tele.Context, state fsm.FSMContext) error {
	g, ok := state.MustGet("g").(scheduleapi.Group)
	if !ok {
		_ = state.Finish(true)
		return c.Send("invalid data, aborting")
	}
	g.Number, _ = strconv.Atoi(c.Data())
	_ = state.Update("g", g)
	state.Set(AcceptGroup)
	markup := keyboards.AcceptMarkup()
	return c.EditOrSend(fmt.Sprintf("Вы выбрали %s.\n\nПодвердите выбор", g.String()), markup)
}

func (h Handler) SGAcceptCallback(c tele.Context, state fsm.FSMContext) error {
	g, ok := state.MustGet("g").(scheduleapi.Group)
	if !ok {
		_ = state.Finish(true)
		return c.Send("invalid data, aborting")
	}
	defer state.Finish(true)

	if err := h.uc.SetGroup(context.TODO(), c.Chat().ID, g); err != nil {
		return c.Send("error: " + err.Error())
	}
	return c.Send("Для данного чата установлена группа: " + g.String())
}

func (h Handler) SGCancelCallback(c tele.Context, state fsm.FSMContext) error {
	_ = state.Finish(true)
	return c.EditOrSend("Выбор группы отменён")
}
