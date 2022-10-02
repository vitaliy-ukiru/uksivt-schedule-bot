package telegram

import (
	"context"

	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram/keyboards"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/groups"
	scheduleapi "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Handler struct {
	uc     chat.Usecase
	uksivt *scheduleapi.Client
	groups groups.Service

	logger *zap.Logger
}

func NewHandler(
	uc chat.Usecase,
	uksivt *scheduleapi.Client,
	groups groups.Groups,
	logger *zap.Logger,
) *Handler {

	return &Handler{
		uc:     uc,
		uksivt: uksivt,
		groups: groups,
		logger: logger,
	}
}

func (h Handler) Route(m *fsm.Manager) {
	m.Use(middleware.AutoRespond())
	m.Group().Handle("/start", h.StartCommand)

	m.Bind("/select_group", fsm.DefaultState, h.SGCommand)
	{
		unique := "\f" + keyboards.SGCallback
		m.Bind(unique, fsm.DefaultState, h.SGYearCallback)
		m.Bind(unique, SelectSpec, h.SGSpecCallback)
		m.Bind(unique, SelectNumber, h.SGNumCallback)
		m.Bind(&keyboards.AcceptBtn, AcceptGroup, h.SGAcceptCallback)
		m.Bind(&keyboards.CancelBtn, fsm.AnyState, h.SGCancelCallback)
	}
	m.Bind("/state", fsm.AnyState, func(c tele.Context, state fsm.FSMContext) error {
		return c.Send(state.State().String())
	})
	m.Bind("/group", fsm.DefaultState, h.GetGroupCommand)
	m.Bind("/lessons", fsm.DefaultState, h.ScheduleCommand)
	keyboards.
		ScheduleCallback.
		MustFilter().
		Handle(m.Group(), h.ScheduleExplorerCallback)
}

const ChatKey = "chat"

func (h Handler) ChatMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		chatObj, err := h.uc.Lookup(context.TODO(), c.Chat().ID)
		if err != nil {
			h.logger.Error("cannot get chat", zap.Error(err))
			return err
		}
		c.Set(ChatKey, chatObj)
		return next(c)
	}
}
