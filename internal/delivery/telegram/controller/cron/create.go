package cron

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/scheduler"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

const SelectTimeText = "Выберите время для отправки: "

var (
	cronSG     = fsm.NewStateGroup("cc")
	SelectOpt  = cronSG.New("opt")
	InputTitle = cronSG.New("name")
	AcceptCron = cronSG.New("accept")
)

type Cron struct {
	At    time.Time
	Title string
	Flags scheduler.FlagSet

	ChatID int64
	Issuer int64
	Step   *fsm.State // step uses for back jumping in steps
}

func (h *Handler) OnlyIssuerMiddleware(m *fsm.Manager) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return m.HandlerAdapter(func(c tele.Context, state fsm.Context) error {
			s := state.State()
			if !fsm.ContainsState(s, cronSG.States...) {
				return next(c)
			}

			cron, ok := state.MustGet("cc").(*Cron)
			if !ok {
				return c.Send("error: fail get data")
			}
			if cron.Issuer != c.Sender().ID {
				return nil // ignore
			}

			return next(c)

		})
	}
}

func (h *Handler) Create(c tele.Context, _ fsm.Context) error {
	markup := AMTimesMarkup(c.Sender().Recipient(), 30*time.Minute)
	return c.Send(SelectTimeText, markup)
}

func (h *Handler) CallbackPM(c tele.Context) error {
	markup := PMTimesMarkup(c.Sender().Recipient(), 30*time.Minute)
	return c.EditOrSend(SelectTimeText, markup)
}

func (h *Handler) CallbackAM(c tele.Context) error {
	markup := AMTimesMarkup(c.Sender().Recipient(), 30*time.Minute)
	return c.Send(SelectTimeText, markup)
}

func (h *Handler) SelectTimeCallback(c tele.Context, state fsm.Context) error {
	m, err := CreateCallback.MustFilter().Process(c)
	if err != nil {
		return answerCallback(c, "cannot parse data: "+err.Error(), true)
	}

	if m["user"] != c.Sender().Recipient() { // ignore if user not issuer
		return nil
	}

	at, err := time.Parse(`15:04`, m["time"])
	if err != nil {
		return c.Send("invalid time data: " + err.Error())
	}

	cron := &Cron{
		ChatID: c.Chat().ID,
		At:     at,
		Issuer: c.Sender().ID,
	}

	if err := state.Update("cc", cron); err != nil {
		return answerCallback(c, "cannot save data: "+err.Error(), true)
	}

	state.Set(SelectOpt)
	cron.Step = &SelectOpt
	return sendFlagsMenu(c, 0)
}

func sendFlagsMenu(c tele.Context, flags scheduler.FlagSet) error {
	const s = `Выберите режим отправки.
Более подробное описание каждого:
	0. Отправлять для следующего дня (совместима с другими режимами)
	1. Отправлять полное расписание
	2. Отправлять только если есть замены
	3. Отправлять если есть замены, если их нет то отправится уведомление об их отсутствии
	4. Отправлять полное расписание только если есть замены`
	return c.EditOrSend(s, FlagsMarkup(flags))
}

func (h *Handler) FlagsCallback(c tele.Context, state fsm.Context) error {
	cron := state.MustGet("cc").(*Cron)

	callback := c.Callback().Data

	selectedFlag, ok := FlagSetFromCallback(callback)
	if !ok {
		return answerCallback(c, "unknown callback: "+callback, true)
	}

	flags := joinFlags(cron.Flags, selectedFlag)

	if err := sendFlagsMenu(c, flags); err != nil {
		return err
	}

	cron.Flags = flags
	return nil
}

func (h *Handler) AcceptFlagsCallback(c tele.Context, state fsm.Context) error {
	cron := state.MustGet("cc").(*Cron)

	if err := c.Edit("Отправьте название данной задачи"); err != nil {
		return err
	}
	state.Set(InputTitle)
	cron.Step = &InputTitle
	return nil
}

func (h *Handler) InputTitle(c tele.Context, state fsm.Context) error {
	cron := state.MustGet("cc").(*Cron)

	title := c.Text()
	cron.Title = title

	err := c.Send(
		fmt.Sprintf(
			"Название: %s\nВремя: %s\nОпции:%s",
			cron.Title,
			cron.At.Format("15:04"),
			flagString(cron.Flags, ";"),
		),
		AcceptMarkup(),
	)
	if err != nil {
		return err
	}
	cron.Step = &AcceptCron
	state.Set(AcceptCron)
	return nil
}

func (h *Handler) AcceptCallback(c tele.Context, state fsm.Context) error {
	cron := state.MustGet("cc").(*Cron)

	chat, _, err := h.chats.Lookup(context.TODO(), c.Chat().ID)
	if err != nil {
		h.logger.Error("cannot get chat", zap.Error(err))
		return c.Send("тех. ошибка")
	}

	dto := cron.ToDTO()
	dto.ChatID = chat.ID

	ctx := context.Background()
	if _, err := h.crons.Create(ctx, dto); err != nil {
		state.Finish(true)
		h.logger.Error("cannot create cron", zap.Error(err))
		return c.Send("Произошла ошибка: " + err.Error())
	}

	return c.Send("Задача создана.")
}

func (c Cron) ToDTO() scheduler.CreateJobDTO {
	return scheduler.CreateJobDTO{
		Title: c.Title,
		At:    c.At,
		Flags: uint16(c.Flags),
	}
}

func joinFlags(flags scheduler.FlagSet, input scheduler.FlagSet) scheduler.FlagSet {
	//TODO: simplify
	if input == scheduler.NextDay {
		if flags.Has(input) {
			return flags.Without(input)
		}

		return flags.With(input)
	}

	result := flags & scheduler.NextDay // save NextDay current state
	if flags.Has(input) {
		return result
	}
	return result.With(input)

}

func flagString(flags scheduler.FlagSet, sep string) string {
	modes := make([]string, 0, 2)
	for _, mode := range FlagModes {
		if flags.Has(mode.Mode) {
			modes = append(modes, mode.Text)
		}
	}
	return strings.Join(modes, sep)
}

func (h *Handler) BackBtnCallback(c tele.Context, state fsm.Context) error {
	cron := state.MustGet("cc").(*Cron)
	switch *cron.Step {
	case SelectOpt:
		cron.Step = nil
		return h.Create(c, state)
	case InputTitle:
		cron.Step = &SelectOpt
		return sendFlagsMenu(c, cron.Flags)
	case AcceptCron:
		return h.AcceptFlagsCallback(c, state)
	}
	return nil
}
