package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	chatPostgres "github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/adapters/storage/chat/postgres"
	cronPostgres "github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/adapters/storage/scheduler/postgres"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/chat"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/config"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/pkg/groups"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/pkg/schedule"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/scheduler"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/client/pg"
	scheduleapi "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

var (
	configPath     = flag.String("config", "./configs/app.yaml", "path to config file")
	groupsFilePath = flag.String("groups", "./configs/groups.json", "path to groups.json")
	envFilePath    = flag.String("env-file", "", "path to .env file")
)

func loadEnv(file string) error {
	if file == "" {
		return nil
	}
	return godotenv.Load(file)
}

func main() {
	flag.Parse()

	{
		if err := loadEnv(*envFilePath); err != nil {
			panic(err)
		}

		fileBytes, err := os.Open(*configPath)
		if err != nil {
			panic(err)
		}

		if err := config.Init(fileBytes); err != nil {
			panic(err)
		}
	}

	cfg := config.Get()

	log, err := newLoggerConfig(cfg).Build()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := log.Sync(); err != nil {
			panic(err)
		}
	}()

	bot, err := tele.NewBot(tele.Settings{
		Token:       cfg.Telegram.Token,
		Poller:      &tele.LongPoller{Timeout: cfg.Telegram.LongPollTimeout},
		Synchronous: false,
		Verbose:     false,
		ParseMode:   tele.ModeHTML,
	})
	if err != nil {
		log.Fatal("cannot init Bot", zap.Error(err))
	}
	log.Info("bot initialized", zap.String("bot_username", bot.Me.Username))

	groupsService, err := groups.NewInMemoryFromFile(*groupsFilePath)
	if err != nil {
		log.Fatal("cannot read groups file", zap.Error(err))
	}

	pool, err := pg.New(context.TODO(), pg.ConnString(
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Database,
		cfg.Database.Host,
		cfg.Database.Port,
	))
	if err != nil {
		log.Fatal("cannot connect to pg pool", zap.Error(err))
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatal("cannot ping pg pool", zap.Error(err))
	}
	{
		poolCfg := pool.Config().ConnConfig
		log.Info(
			"connected to postgres",
			zap.String("host", poolCfg.Host),
			zap.Uint16("port", poolCfg.Port),
			zap.String("user", poolCfg.User),
		)

	}

	var (
		chatStorage = chatPostgres.NewRepository(pool)
		chatService = chat.NewService(chatStorage)
	)
	log.Debug("chat services configured")

	var (
		uksivtClient   = scheduleapi.NewClient(nil, cfg.UKSIVT.ApiURL)
		uksivtSchedule = schedule.NewService(uksivtClient)
	)
	log.Debug("uksivt services configured")

	var (
		cronStorage = cronPostgres.NewRepository(pool)
		cronService = scheduler.NewService(cronStorage, cfg.Scheduler.Range)
	)
	log.Debug("cron services configured")

	handler := telegram.NewHandler(
		chatService,
		uksivtSchedule,
		groupsService,
		cronService,
		cfg,
		log,
		bot,
	)

	location, err := time.LoadLocation(cfg.Scheduler.TimeLocation)
	if err != nil {
		log.Fatal("cannot parse location", zap.Error(err))
	}

	cron := gocron.NewScheduler(location)
	storage := memory.NewStorage()
	m := fsm.NewManager(bot.Group(), storage)

	{
		handler.Route(m)
		if err := handler.Schedule(cron); err != nil {
			log.Fatal("cannot schedule task", zap.Error(err))
		}
	}

	{
		log.Info("start listening")
		cron.StartAsync()
		bot.Start()
	}
}
