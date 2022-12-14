package telegram

import (
	"github.com/go-co-op/gocron"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/chat"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/config"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram/controller/cron"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram/controller/group"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram/controller/schedule"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Handler struct {
	uc chat.Usecase

	groups  *group.Handler
	crons   *cron.Handler
	lessons *schedule.Handler

	cfg    *config.Config
	logger *zap.Logger
	bot    *tele.Bot
}

func New(
	uc chat.Usecase,
	groups *group.Handler,
	crons *cron.Handler,
	lessons *schedule.Handler,
	cfg *config.Config,
	logger *zap.Logger,
	bot *tele.Bot,
) *Handler {
	return &Handler{
		uc:      uc,
		groups:  groups,
		crons:   crons,
		lessons: lessons,
		cfg:     cfg,
		logger:  logger,
		bot:     bot,
	}
}

func (h Handler) BindHandlers(m *fsm.Manager) {
	m.Use(middleware.AutoRespond())
	m.Group().Handle("/start", h.StartCommand)

	m.Bind("/state", fsm.AnyState, func(c tele.Context, state fsm.Context) error {
		return c.Send(state.State().String())
	})

	m.Group().Handle("/cron", func(c tele.Context) error {
		if c.Sender().ID == h.cfg.Telegram.AdminID {
			h.lessons.CronSchedulerJob()
		}
		return nil
	})

	{
		h.lessons.Bind(m.NewGroup())
		h.crons.Bind(m.NewGroup())
		h.groups.Bind(m.NewGroup())
	}
}

func (h Handler) BindCrons(s *gocron.Scheduler) error {
	_, err := s.Cron(h.cfg.Scheduler.Cron).Do(h.lessons.CronSchedulerJob)
	return err
}
