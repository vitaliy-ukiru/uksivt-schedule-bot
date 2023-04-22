package group

import (
	"context"
	"fmt"
	"strconv"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/group"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

var (
	groupSG = fsm.NewStateGroup("s.g")

	SelectSpecState   = groupSG.New("spec")
	SelectNumberState = groupSG.New("num")
	AcceptGroupState  = groupSG.New("accept")
)

func (h *Handler) Command(c tele.Context, _ fsm.Context) error {
	markup, err := h.getYearsMarkup(context.TODO())
	if err != nil {
		return c.Send("ERROR: cannot get groups(year): " + err.Error())
	}
	return c.Send("Выберите год поступления:", markup)
}

func (h *Handler) YearCallback(c tele.Context, state fsm.Context) error {
	yearStr := c.Data()
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		h.logger.Error("invalid  year callback", zap.Error(err))
		return c.Send("ERROR: parse year: " + err.Error())
	}
	_ = state.Update("g", group.Group{Year: year})

	specsMarkup, err := h.getSpecsMarkup(context.TODO(), year)
	if err != nil {
		return c.Send("ERROR: cannot get groups(spec): " + err.Error())
	}
	state.Set(SelectSpecState)
	return c.EditOrSend("Выберите специальность", specsMarkup)
}

func (h *Handler) SpecCallback(c tele.Context, state fsm.Context) error {
	g, ok := state.MustGet("g").(group.Group)
	if !ok {
		_ = state.Finish(true)
		return c.Send("ERROR: invalid data, aborting")
	}
	g.Spec = c.Data()
	_ = state.Update("g", g)

	numsMarkup, err := h.getNumsMarkup(context.TODO(), g)
	if err != nil {
		return c.Send("ERROR: cannot get groups(num): " + err.Error())
	}
	state.Set(SelectNumberState)
	return c.EditOrSend("Выберите группу", numsMarkup)
}

func (h *Handler) NumCallback(c tele.Context, state fsm.Context) error {
	g, ok := state.MustGet("g").(group.Group)
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
	g, ok := state.MustGet("g").(group.Group)
	if !ok {
		_ = state.Finish(true)
		return c.Send("invalid data, aborting")
	}
	defer state.Finish(true)

	if err := h.chats.SetGroup(context.Background(), c.Chat().ID, g.String()); err != nil {
		return c.Send("error: " + err.Error())
	}
	return c.Send("Для данного чата установлена группа: " + g.String())
}

func (h *Handler) CancelCallback(c tele.Context, state fsm.Context) error {
	_ = state.Finish(true)
	return c.EditOrSend("Выбор группы отменён")
}

// BackToYearsCallback returns to input year menu
func (h *Handler) BackToYearsCallback(c tele.Context, state fsm.Context) error {
	markup, err := h.getYearsMarkup(context.TODO())
	if err != nil {
		return c.Send("ERROR: cannot get groups(year): " + err.Error())
	}
	state.Finish(false)
	return c.EditOrSend("Выберите год поступления:", markup)
}

// BackToSpecsCallback returns to select spec menu
func (h *Handler) BackToSpecsCallback(c tele.Context, state fsm.Context) error {
	g, ok := state.MustGet("g").(group.Group)
	if !ok {
		_ = state.Finish(true)
		return c.Send("ERROR: invalid data, aborting")
	}

	specsMarkup, err := h.getSpecsMarkup(context.TODO(), g.Year)
	if err != nil {
		return c.Send("ERROR: cannot get groups(spec): " + err.Error())
	}
	state.Set(SelectSpecState)
	return c.EditOrSend("Выберите специальность", specsMarkup)
}

// BackToNumbersCallback returns to select group number menu
func (h *Handler) BackToNumbersCallback(c tele.Context, state fsm.Context) error {
	g, ok := state.MustGet("g").(group.Group)
	if !ok {
		_ = state.Finish(true)
		return c.Send("ERROR: invalid data, aborting")
	}

	numsMarkup, err := h.getNumsMarkup(context.TODO(), g)
	if err != nil {
		return c.Send("ERROR: cannot get groups(num): " + err.Error())
	}
	state.Set(SelectNumberState)
	return c.EditOrSend("Выберите группу", numsMarkup)
}

func (h *Handler) getYearsMarkup(ctx context.Context) (*tele.ReplyMarkup, error) {
	years, err := h.groups.Years(ctx)
	if err != nil {
		h.logger.Error("get years", zap.Error(err))
		return nil, err
	}
	markup := SelectYearMarkup(years)
	return markup, nil
}

func (h *Handler) getSpecsMarkup(ctx context.Context, year int) (*tele.ReplyMarkup, error) {
	specs, err := h.groups.Specs(ctx, year)
	if err != nil {
		h.logger.Error("get specs", zap.Error(err), zap.Int("year", year))
		return nil, err
	}

	markup := SelectSpecMarkup(year, specs)
	return markup, nil
}

func (h *Handler) getNumsMarkup(ctx context.Context, g group.Group) (*tele.ReplyMarkup, error) {
	numbers, err := h.groups.Numbers(ctx, g.Year, g.Spec)
	if err != nil {
		h.logger.Error("get nums",
			zap.Error(err),
			zap.Int("year", g.Year),
			zap.String("spec", g.Spec),
		)
		return nil, err
	}
	markup := SelectNumberMarkup(g.Year, g.Spec, numbers)
	return markup, nil
}
