package cron

import (
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/chat"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/scheduler"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type CreateCronHandler struct {
	chats  chat.Usecase
	crons  scheduler.Usecase
	logger *zap.Logger
}

func NewCreateHandler(chats chat.Usecase, crons scheduler.Usecase, logger *zap.Logger) *CreateCronHandler {
	return &CreateCronHandler{chats: chats, crons: crons, logger: logger}
}

func (h *CreateCronHandler) Bind(m *fsm.Manager) {
	m.Bind("/create", fsm.DefaultState, h.CreateCommand)
	m.Group().Handle("/crons", h.ListCrons)

	m.Group().Handle(&PMButton, h.PMCallback)
	m.Group().Handle(&AMButton, h.AMCallback)

	m.Bind(SelectTimeCallback, fsm.DefaultState, h.SelectTimeCallback)
	m.Bind(FlagsCallback, SelectOpt, h.FlagsCallback, h.OnlyIssuerMiddleware(m))

	m.Bind(&AcceptFlags, SelectOpt, h.AcceptFlagsCallback, h.OnlyIssuerMiddleware(m))
	m.Bind(tele.OnText, InputTitle, h.InputTitle, h.OnlyIssuerMiddleware(m))
	m.Bind(&AcceptBtn, AcceptCron, h.AcceptCallback, h.OnlyIssuerMiddleware(m))
}

func (h *EditCronHandler) Bind(m *fsm.Manager) {
	m.Bind("/edit", fsm.DefaultState, h.EditCommand)
	m.Bind(SelectCronEditCallback, fsm.DefaultState, h.EditSelectCronCallback)

	m.Bind(&SelectEditTitle, SelectEditingField, h.EditTitleCallback)
	m.Bind(&SelectEditTime, SelectEditingField, h.EditTimeCallback)
	m.Bind(&SelectEditFlags, SelectEditingField, h.EditFlagsCallback)
	m.Bind(&DoneEditing, SelectEditingField, h.DoneEditingCallback)

	m.Bind(tele.OnText, EditTitle, h.InputNewTitle)
	m.Bind(SelectTimeCallback, EditTime, h.InputTime)
	m.Bind(FlagsCallback, EditFlags, h.InputFlagCallback)
	m.Bind(&AcceptFlags, EditFlags, h.AcceptFlagCallback)

}

func answerCallback(c tele.Context, text string, alert bool) error {
	return c.Respond(&tele.CallbackResponse{
		Text:      text,
		ShowAlert: alert,
	})
}
