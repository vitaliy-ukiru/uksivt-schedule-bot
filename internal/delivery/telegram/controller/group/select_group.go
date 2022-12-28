package group

import (
	"context"
	"fmt"
	"strconv"

	"github.com/vitaliy-ukiru/fsm-telebot"
	scheduleapi "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
	tele "gopkg.in/telebot.v3"
)

var (
	//SelectYear   = groupSG.New("year")
	groupSG = fsm.NewStateGroup("s.g")

	SelectSpecState   = groupSG.New("spec")
	SelectNumberState = groupSG.New("num")
	AcceptGroupState  = groupSG.New("accept")
)

func (h *Handler) Command(c tele.Context, _ fsm.Context) error {
	//chat, err := h.uc.ByTelegramID(context.TODO(), c.Chat().ID)
	//if err != nil {
	//	return c.Send("cannot get chat: " + err.Error())
	//}
	//state.Set(SelectYear)
	markup := SelectYearMarkup(h.groups.Years())
	return c.Send("Выберите год поступления:", markup)
}

func (h *Handler) YearCallback(c tele.Context, state fsm.Context) error {
	yearStr := c.Data()
	year, _ := strconv.Atoi(yearStr)
	//if err != nil {
	//	return err
	//}
	state.Set(SelectSpecState)
	_ = state.Update("g", scheduleapi.Group{Year: year})
	specs := h.groups.Specs(year)
	markup := SelectSpecMarkup(year, specs)

	return c.EditOrSend("Выберите специальность", markup)
}

func (h *Handler) SpecCallback(c tele.Context, state fsm.Context) error {
	//year, ok := state.MustGet("year").(int)
	g, ok := state.MustGet("g").(scheduleapi.Group)
	if !ok {
		_ = state.Finish(true)
		return c.Send("invalid data, aborting")
	}
	g.Spec = c.Data()

	_ = state.Update("g", g)
	state.Set(SelectNumberState)
	numbers := h.groups.Numbers(g.Year, g.Spec)
	markup := SelectNumberMarkup(g.Year, g.Spec, numbers)
	return c.EditOrSend("Выберите группу", markup)
}

func (h *Handler) NumCallback(c tele.Context, state fsm.Context) error {
	g, ok := state.MustGet("g").(scheduleapi.Group)
	if !ok {
		_ = state.Finish(true)
		return c.Send("invalid data, aborting")
	}
	g.Number, _ = strconv.Atoi(c.Data())
	_ = state.Update("g", g)
	state.Set(AcceptGroupState)
	markup := AcceptMarkup()
	return c.EditOrSend(fmt.Sprintf("Вы выбрали %s.\n\nПодвердите выбор", g.String()), markup)
}

func (h *Handler) AcceptCallback(c tele.Context, state fsm.Context) error {
	g, ok := state.MustGet("g").(scheduleapi.Group)
	if !ok {
		_ = state.Finish(true)
		return c.Send("invalid data, aborting")
	}
	defer state.Finish(true)

	if err := h.chats.SetGroup(context.Background(), c.Chat().ID, g); err != nil {
		return c.Send("error: " + err.Error())
	}
	return c.Send("Для данного чата установлена группа: " + g.String())
}

func (h *Handler) CancelCallback(c tele.Context, state fsm.Context) error {
	_ = state.Finish(true)
	return c.EditOrSend("Выбор группы отменён")
}
