package cron

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/config"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/scheduler"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

const SelectTimeText = "Выберите время для отправки: "

var (
	createSG   = fsm.NewStateGroup("cc")
	SelectOpt  = createSG.New("opt")
	InputTitle = createSG.New("name")
	AcceptCron = createSG.New("accept")
)

type Cron struct {
	At    time.Time
	Title string
	Flags scheduler.FlagSet

	ChatID int64
	Issuer int64
	Step   *fsm.State // step uses for back jumping in steps
}

func (h *CreateCronHandler) OnlyIssuerMiddleware(m *fsm.Manager) tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return m.HandlerAdapter(func(c tele.Context, state fsm.Context) error {
			s := state.State()
			if !fsm.ContainsState(s, createSG.States...) {
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

func (h *CreateCronHandler) CreateCommand(c tele.Context, _ fsm.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	chat, _, err := h.chats.Lookup(ctx, c.Chat().ID)
	if err != nil {
		h.logger.Error("cannot get chat", zap.Error(err))
		return c.Send("Не могу получить информацию о чате. " +
			"Попробуйте позже или сообщите разработчику (/help)")
	}

	if chat.Group == nil {
		return c.Send("Нельзя создать программу пока не установлена группа для чата (/select_group)")
	}

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cronsCount, err := h.crons.CountInChat(ctx, chat.ID)
	if err != nil {
		h.logger.Error("cannot get count of crons", zap.Error(err))
		return c.Send("Произошла ошибка. Попробуйте позже")
	}

	if cronsCount >= config.MaxCrons {
		return c.Send("Нельзя создавать больше программ. Превышен лимит.")
	}

	markup := TimesMarkupAM(c.Sender().Recipient(), 30*time.Minute)
	return c.Send(SelectTimeText, markup)
}

func (h *CreateCronHandler) PMCallback(c tele.Context) error {
	markup := TimesMarkupPM(c.Sender().Recipient(), 30*time.Minute)
	return c.EditOrSend(SelectTimeText, markup)
}

func (h *CreateCronHandler) AMCallback(c tele.Context) error {
	markup := TimesMarkupAM(c.Sender().Recipient(), 30*time.Minute)
	return c.EditOrSend(SelectTimeText, markup)
}

func (h *CreateCronHandler) SelectTimeCallback(c tele.Context, state fsm.Context) error {
	m, err := SelectTimeCallback.MustFilter().Process(c)
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

	cron := Cron{
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

func (h *CreateCronHandler) FlagsCallback(c tele.Context, state fsm.Context) error {
	cron := state.MustGet("cc").(Cron)

	callback := c.Callback().Data
	selectedFlag, ok := FlagSetFromCallback(callback)
	if !ok {
		return answerCallback(c, "unknown callback: "+callback, true)
	}

	flags := joinFlags(cron.Flags, selectedFlag)
	cron.Flags = flags
	state.Update("cc", cron)

	return sendFlagsMenu(c, flags)
}

func (h *CreateCronHandler) AcceptFlagsCallback(c tele.Context, state fsm.Context) error {
	cron := state.MustGet("cc").(Cron)
	defer state.Update("cc", cron)

	state.Set(InputTitle)
	cron.Step = &InputTitle

	return c.EditOrSend("Отправьте название данной задачи")
}

func (h *CreateCronHandler) InputTitle(c tele.Context, state fsm.Context) error {
	cron := state.MustGet("cc").(Cron)

	cron.Title = c.Text()
	cron.Step = &AcceptCron
	state.Update("cc", cron)
	state.Set(AcceptCron)

	return c.Send(
		fmt.Sprintf(
			"Название: %s\nВремя: %s\nОпции:%s",
			cron.Title,
			cron.At.Format("15:04"),
			flagString(cron.Flags, ";"),
		),
		AcceptMarkup(),
	)
}

func (h *CreateCronHandler) AcceptCallback(c tele.Context, state fsm.Context) error {
	go c.Respond()
	cron := state.MustGet("cc").(Cron)

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
	state.Finish(true)
	return c.Send("Задача создана.")
}

func (h *CreateCronHandler) CancelCreateCallback(c tele.Context, state fsm.Context) error {
	go c.Respond()
	state.Finish(true)
	return c.EditOrSend("Создание отменено")
}

func (c Cron) ToDTO() scheduler.CreateJobDTO {
	return scheduler.CreateJobDTO{
		Title: c.Title,
		At:    c.At,
		Flags: uint16(c.Flags),
	}
}

func joinFlags(flags scheduler.FlagSet, input scheduler.FlagSet) scheduler.FlagSet {
	if input == scheduler.NextDay {
		return flags.Toggle(input)
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

func (h *CreateCronHandler) BackBtnCallback(c tele.Context, state fsm.Context) error {
	cron := state.MustGet("cc").(Cron)
	defer state.Update("cc", cron)
	switch *cron.Step {
	case SelectOpt:
		cron.Step = nil
		return h.CreateCommand(c, state)
	case InputTitle:
		cron.Step = &SelectOpt
		return sendFlagsMenu(c, cron.Flags)
	case AcceptCron:
		return h.AcceptFlagsCallback(c, state)
	}
	return nil
}
