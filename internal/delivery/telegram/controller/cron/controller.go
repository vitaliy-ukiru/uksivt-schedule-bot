package cron

import (
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/chat"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/scheduler"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type Handler struct {
	chats  chat.Usecase
	crons  scheduler.Usecase
	logger *zap.Logger
}

func New(chats chat.Usecase, crons scheduler.Usecase, logger *zap.Logger) *Handler {
	return &Handler{chats: chats, crons: crons, logger: logger}
}

func (h *Handler) Bind(m *fsm.Manager) {
	m.Bind("/create", fsm.DefaultState, h.Create)
	m.Group().Handle("/crons", h.ListCrons)

	m.Group().Handle(&PMButton, h.CallbackPM)
	m.Group().Handle(&AMButton, h.CallbackAM)

	m.Bind(CreateCallback, fsm.DefaultState, h.SelectTimeCallback)
	m.Bind(FlagsCallback, SelectOpt, h.FlagsCallback, h.OnlyIssuerMiddleware(m))

	m.Bind(&AcceptFlags, SelectOpt, h.AcceptFlagsCallback, h.OnlyIssuerMiddleware(m))
	m.Bind(tele.OnText, InputTitle, h.InputTitle, h.OnlyIssuerMiddleware(m))
	m.Bind(&AcceptBtn, AcceptCron, h.AcceptCallback, h.OnlyIssuerMiddleware(m))
}

func answerCallback(c tele.Context, text string, alert bool) error {
	return c.Respond(&tele.CallbackResponse{
		Text:      text,
		ShowAlert: alert,
	})
}
