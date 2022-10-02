package main

import (
	"context"
	"flag"
	"os"
	"time"

	"github.com/joho/godotenv"
	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/config"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat/storage/postgres"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/groups"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/client/pg"
	scheduleapi "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

var (
	configPath     = flag.String("config", "./config/app.yaml", "path to config file")
	envFilePath    = flag.String("env-file", "", "path to .env file")
	groupsFilePath = flag.String("groups", "./data/groups.json", "path to groups.json")
)

func loadEnv() error {
	if *envFilePath != "" {
		return godotenv.Load(*envFilePath)
	}
	return nil
}

func main() {
	flag.Parse()
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	if err := loadEnv(); err != nil {
		log.Fatal("cannot read env file",
			zap.String("path", *envFilePath), zap.Error(err),
		)
	}

	fileBytes, err := os.Open(*configPath)
	if err != nil {
		log.Fatal("cannot read config file",
			zap.String("path", *configPath), zap.Error(err),
		)
	}

	if err := config.Init(fileBytes); err != nil {
		log.Fatal("cannot init config", zap.Error(err))
	}
	cfg := config.Get()

	bot, err := tele.NewBot(tele.Settings{
		Token:       cfg.Telegram.Token,
		Poller:      &tele.LongPoller{Timeout: 2 * time.Second},
		Synchronous: false,
		Verbose:     true,
	})
	if err != nil {
		log.Fatal("cannot init Bot", zap.Error(err))
	}

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

	chatStorage := postgres.NewRepository(pool)
	chatService := chat.NewService(chatStorage)
	uksivtSchedule := scheduleapi.NewClient(nil, cfg.Schedule.ApiURL)
	handler := telegram.NewHandler(chatService, uksivtSchedule, groupsService, log)

	storage := memory.NewStorage()
	m := fsm.NewManager(bot.Group(), storage)
	handler.Route(m)
	bot.Start()
}
