package schedule

import (
	"context"
	"time"

	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/adapters/schedule"
	chat2 "github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/scheduler"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

type Handler struct {
	chats  chat2.Usecase
	uksivt schedule.Usecase
	crons  scheduler.Usecase

	logger *zap.Logger
	bot    *tele.Bot
}

func New(
	chats chat2.Usecase,
	uksivt schedule.Usecase,
	crons scheduler.Usecase,
	logger *zap.Logger,
	bot *tele.Bot,
) *Handler {
	return &Handler{
		chats:  chats,
		uksivt: uksivt,
		crons:  crons,
		logger: logger,
		bot:    bot,
	}
}

func (h *Handler) Bind(m *fsm.Manager) {
	m.Bind("/lessons", fsm.DefaultState, h.LessonsCommand)
	Callback.MustFilter().Handle(
		m.Group(),
		h.ExplorerCallback,
	)
}

func (h *Handler) getChat(tgID int64) *chat2.Chat {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	chatObj, _, err := h.chats.Lookup(ctx, tgID)
	if err != nil {
		h.logger.Error("cannot get chat", zap.Error(err))
		return nil
	}

	return chatObj
}
