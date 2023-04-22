package group

import (
	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/adapters/group"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat"
	"go.uber.org/zap"
)

type Handler struct {
	chats  chat.Usecase
	groups group.Fetcher

	logger *zap.Logger
}

func New(chats chat.Usecase, groups group.Fetcher, logger *zap.Logger) *Handler {
	return &Handler{chats: chats, groups: groups, logger: logger}
}

func (h *Handler) Bind(m *fsm.Manager) {
	m.Group().Handle("/group", h.GetGroupCommand)

	m.Bind("/select_group", fsm.DefaultState, h.Command)
	{
		unique := "\f" + SGCallback
		m.Bind(unique, fsm.DefaultState, h.YearCallback)
		m.Bind(unique, SelectSpecState, h.SpecCallback)
		m.Bind(unique, SelectNumberState, h.NumCallback)
		m.Bind(&AcceptBtn, AcceptGroupState, h.AcceptCallback)
		m.Bind(&CancelBtn, fsm.AnyState, h.CancelCallback)

		{
			backBtn := &BackBtn
			m.Bind(backBtn, SelectSpecState, h.BackToYearsCallback)
			m.Bind(backBtn, SelectNumberState, h.BackToSpecsCallback)
			m.Bind(backBtn, AcceptGroupState, h.BackToNumbersCallback)
		}
	}

}
