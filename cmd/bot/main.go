package main

import (
	"context"
	"encoding/json"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	chatPostgres "github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/adapters/dao/chat/postgres"
	groupPostgres "github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/adapters/dao/group/postgres"
	cronPostgres "github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/adapters/dao/scheduler/postgres"
	groupAdapter "github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/adapters/group"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/adapters/schedule"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/config"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/group"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/scheduler"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/client/pg"
	scheduleapi "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

var (
	configPath     = flag.String("config", "./configs/app.yaml", "path to config file")
	envFilePath    = flag.String("env-file", "", "path to .env file")
	groupsFilePath = flag.String("groups", "./configs/groups.json",
		"path to groups.json. Use \"postgres\" for using PostgresSQL (before execute db/migration/groups.up.sql)")
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
		_ = log.Sync()
	}()

	bot, err := tele.NewBot(tele.Settings{
		Token:       cfg.Telegram.Token,
		Poller:      &tele.LongPoller{Timeout: cfg.Telegram.LongPollTimeout},
		Synchronous: false,
		Verbose:     false,
		ParseMode:   tele.ModeHTML,
		OnError: func(err error, c tele.Context) {
			log := log.With(zap.Error(err), zap.Int64("chat_id", c.Chat().ID))

			//. //Error("error in handler")

			if handlerName := c.Get("handler"); handlerName != nil {
				log = log.With(zap.Any("handler", handlerName))
			}
			log.Warn("in handler error")
			update, errJson := json.Marshal(c.Update())
			if errJson == nil {
				log.Debug(string(update))
			}
		},
	})
	if err != nil {
		log.Fatal("cannot init Bot", zap.Error(err))
	}
	log.Info("bot initialized", zap.String("bot_username", bot.Me.Username))

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
		groupStorage = groupPostgres.NewRepository(pool)
		groupService = group.NewService(groupStorage)
		groupFetcher groupAdapter.Fetcher
	)

	if *groupsFilePath == "postgres" {
		groupFetcher = groupService
	} else {
		groupFetcher, err = groupAdapter.NewInMemoryFromFile(*groupsFilePath)
		if err != nil {
			log.Fatal("cannot read groups file", zap.Error(err))
		}
	}

	log.Debug("group services configured")

	var (
		chatStorage = chatPostgres.NewRepository(pool)
		chatService = chat.NewService(chatStorage, groupService, log)
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

	var (
		groupsHandler     = telegram.NewGroupHandler(chatService, groupFetcher, log)
		cronCreateHandler = telegram.NewCreateCronHandler(chatService, cronService, log)
		cronEditHandler   = telegram.NewEditCronHandler(chatService, cronService, log)
		lessonsHandler    = telegram.NewScheduleHandler(chatService, uksivtSchedule, cronService, log, bot)
	)

	handler := telegram.New(
		chatService,
		groupsHandler,
		cronCreateHandler,
		cronEditHandler,
		lessonsHandler,
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
	m := fsm.NewManager(bot, nil, storage, nil)

	{
		handler.BindHandlers(m)
		if err := handler.Schedule(cron); err != nil {
			log.Fatal("cannot schedule task", zap.Error(err))
		}
	}

	{
		log.Info("start listening")
		cron.StartAsync()
		go bot.Start()
	}

	{
		quit := make(chan os.Signal, 1)
		signal.Notify(quit,
			syscall.SIGINT,
			syscall.SIGQUIT,
			syscall.SIGKILL,
			syscall.SIGTERM,
		)
		sig := <-quit
		log.Warn("shutdown app", zap.String("signal", sig.String()))
	}

	{
		cron.Stop()
		bot.Stop()
		log.Info("stopped bot")
	}
}
