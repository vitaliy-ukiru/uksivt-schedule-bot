package cron

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/scheduler"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type EditCronHandler struct {
	chats  chat.Usecase
	crons  scheduler.Usecase
	logger *zap.Logger
}

func NewEditHandler(chats chat.Usecase, crons scheduler.Usecase, logger *zap.Logger) *EditCronHandler {
	return &EditCronHandler{chats: chats, crons: crons, logger: logger}
}

var (
	editSg             = fsm.NewStateGroup("edit_c")
	SelectEditingField = editSg.New("select_field")
	EditTitle          = editSg.New("input_title")
	EditTime           = editSg.New("input_time")
	EditFlags          = editSg.New("input_flags")
)

func (h *EditCronHandler) EditCommand(c tele.Context, _ fsm.Context) error {
	ctx := context.Background()
	chatObj, _, err := h.chats.Lookup(ctx, c.Chat().ID)
	if err != nil {
		return c.Send(err.Error())
	}

	crons, err := h.crons.ForChat(ctx, chatObj.ID)
	if err != nil {
		return c.Send("cannot get crons: " + err.Error())
	}
	return c.Send("Выберите задачу", SelectCronsMarkup(crons))
}

func (h *EditCronHandler) EditSelectCronCallback(c tele.Context, state fsm.Context) error {
	cronId, err := strconv.ParseInt(c.Callback().Data, 10, 64)
	if err != nil {
		return answerCallback(c, "invalid callback cron id: "+err.Error(), true)
	}

	ctx := context.TODO()
	cron, err := h.crons.ByID(ctx, cronId)
	if err != nil {
		return answerCallback(c, "cannot get cron: "+err.Error(), true)
	}

	state.Set(SelectEditingField)
	state.Update("ce", *cron)
	return h.SendCronView(c, *cron)
}

func (h *EditCronHandler) SendCronView(c tele.Context, cron scheduler.CronJob) error {
	return c.EditOrSend(fmt.Sprintf(
		"Название: %s\nВремя:%s\nОпции:%s\n\nВыберите что редактировать:",
		cron.Title,
		cron.At.Format("15:04"),
		flagString(cron.Flags, "; "),
	), SelectEditingFieldMarkup)
}

func (h *EditCronHandler) EditTitleCallback(c tele.Context, state fsm.Context) error {
	state.Set(EditTitle)
	return c.EditOrSend("Введите новое название")
}

func (h *EditCronHandler) EditTimeCallback(c tele.Context, state fsm.Context) error {
	state.Set(EditTime)
	return c.EditOrSend(SelectTimeText, TimesMarkupAM(c.Sender().Recipient(), 30*time.Minute))
}

func (h *EditCronHandler) EditFlagsCallback(c tele.Context, state fsm.Context) error {
	state.Set(EditFlags)

	var cron scheduler.CronJob
	if err := state.Get("ce", &cron); err != nil {
		_ = state.Finish(true)
		return answerCallback(c, "cannot get cron from context: "+err.Error(), true)
	}

	state.Update("fce", cron.Flags)
	return sendFlagsMenu(c, cron.Flags)

}

func (h *EditCronHandler) InputNewTitle(c tele.Context, state fsm.Context) error {
	var cron scheduler.CronJob
	if err := state.Get("ce", &cron); err != nil {
		_ = state.Finish(true)
		return answerCallback(c, "cannot get cron from context: "+err.Error(), true)
	}
	cron.Title = c.Text()
	state.Update("ce", cron)
	state.Set(SelectEditingField)

	return h.SendCronView(c, cron)
}

func (h *EditCronHandler) InputTime(c tele.Context, state fsm.Context) error {
	m, err := SelectTimeCallback.MustFilter().Process(c)
	if err != nil {
		_ = state.Finish(true)
		return answerCallback(c, "cannot parse data: "+err.Error(), true)
	}

	if m["user"] != c.Sender().Recipient() { // ignore if user not issuer
		return nil
	}

	at, err := time.Parse(`15:04`, m["time"])
	if err != nil {
		return c.Send("invalid time data: " + err.Error())
	}

	var cron scheduler.CronJob
	if err := state.Get("ce", &cron); err != nil {
		_ = state.Finish(true)
		return answerCallback(c, "cannot get cron from context: "+err.Error(), true)
	}

	cron.At = at
	state.Update("ce", cron)
	state.Set(SelectEditingField)

	return h.SendCronView(c, cron)
}

func (h *EditCronHandler) InputFlagCallback(c tele.Context, state fsm.Context) error {
	var flags scheduler.FlagSet
	if err := state.Get("fce", &flags); err != nil {
		_ = state.Finish(true)
		return answerCallback(c, "cannot get cron from context: "+err.Error(), true)
	}
	callback := c.Callback().Data

	selectedFlag, ok := FlagSetFromCallback(callback)
	if !ok {
		return answerCallback(c, "unknown callback: "+callback, true)
	}

	flags = joinFlags(flags, selectedFlag)
	state.Update("fce", flags)

	return sendFlagsMenu(c, flags)
}

func (h *EditCronHandler) AcceptFlagCallback(c tele.Context, state fsm.Context) error {
	var cron scheduler.CronJob
	if err := state.Get("ce", &cron); err != nil {
		_ = state.Finish(true)
		return answerCallback(c, "cannot get cron from context: "+err.Error(), true)
	}

	var flags scheduler.FlagSet
	if err := state.Get("fce", &flags); err != nil {
		_ = state.Finish(true)
		return answerCallback(c, "cannot get cron from context: "+err.Error(), true)
	}

	if flags <= scheduler.NextDay {
		flags = flags.With(scheduler.Full)
	}

	cron.Flags = flags
	state.Update("fce", nil)
	state.Update("ce", cron)

	state.Set(SelectEditingField)
	return h.SendCronView(c, cron)
}

func (h *EditCronHandler) DoneEditingCallback(c tele.Context, state fsm.Context) error {
	var cron scheduler.CronJob
	if err := state.Get("ce", &cron); err != nil {
		_ = state.Finish(true)
		return answerCallback(c, "cannot get cron from context: "+err.Error(), true)
	}

	state.Finish(true)
	ctx := context.Background()
	if err := h.crons.Update(ctx, cron); err != nil {
		return c.EditOrSend("не получилось сохранить: " + err.Error())
	}

	return c.EditOrSend("Изменения сохранены.")
}

func (h *EditCronHandler) CancelEditingCallback(c tele.Context, state fsm.Context) error {
	_ = state.Finish(true)
	return c.EditOrSend("Изменение отменено")
}
