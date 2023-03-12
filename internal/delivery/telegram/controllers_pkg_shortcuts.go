package telegram

import (
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/adapters/group"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/adapters/schedule"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram/controller/cron"
	selectGroup "github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram/controller/group"
	lessons "github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram/controller/schedule"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/scheduler"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func NewCreateCronHandler(chats chat.Usecase, crons scheduler.Usecase, logger *zap.Logger) *cron.CreateCronHandler {
	return cron.NewCreateHandler(chats, crons, logger)
}

func NewEditCronHandler(chats chat.Usecase, crons scheduler.Usecase, logger *zap.Logger) *cron.EditCronHandler {
	return cron.NewEditHandler(chats, crons, logger)
}
func NewGroupHandler(chats chat.Usecase, groups group.Fetcher, logger *zap.Logger) *selectGroup.Handler {
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
