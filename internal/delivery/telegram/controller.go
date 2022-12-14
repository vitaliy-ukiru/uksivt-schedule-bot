package telegram

import (
	"github.com/go-co-op/gocron"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/chat"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/config"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram/keyboards"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/pkg/groups"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/pkg/schedule"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/scheduler"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
)

type Handler struct {
	uc     chat.Usecase
	uksivt schedule.Usecase
	groups groups.Service
	crons  scheduler.CronFetcher

	cfg    *config.Config
	logger *zap.Logger
	bot    *tele.Bot
}

func NewHandler(
	uc chat.Usecase,
	uksivt schedule.Usecase,
	groups groups.Service,
	crons scheduler.CronFetcher,
	cfg *config.Config,
	logger *zap.Logger,
	bot *tele.Bot,
) *Handler {
	return &Handler{
		uc:     uc,
		uksivt: uksivt,
		groups: groups,
		cfg:    cfg,
		logger: logger,
		bot:    bot,
		crons:  crons,
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
		ScheduleCallback.MustFilter().
		Handle(m.Group(), h.ScheduleExplorerCallback)
}

func (h Handler) Schedule(s *gocron.Scheduler) error {
	_, err := s.Cron(h.cfg.Scheduler.Cron).Do(h.CronJobSchedule)
	return err
}
