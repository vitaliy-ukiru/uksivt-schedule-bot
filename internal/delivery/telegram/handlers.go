package telegram

import (
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/adapters/group"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/adapters/schedule"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/chat"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram/controller/cron"
	selectGroup "github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram/controller/group"
	lessons "github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram/controller/schedule"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/scheduler"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func NewCronHandler(chats chat.Usecase, crons scheduler.Usecase, logger *zap.Logger) *cron.Handler {
	return cron.New(chats, crons, logger)
}

func NewGroupHandler(chats chat.Usecase, groups group.Usecase, logger *zap.Logger) *selectGroup.Handler {
	return selectGroup.New(chats, groups, logger)
}

func NewScheduleHandler(
	chats chat.Usecase,
	uksivt schedule.Usecase,
	crons scheduler.Usecase,
	logger *zap.Logger,
	bot *tele.Bot,
) *lessons.Handler {
	return lessons.New(chats, uksivt, crons, logger, bot)
}
