package main

import (
	"context"
	"time"

	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/fsm-telebot/storages/memory"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat/storage/postgres"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/groups"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/client/pg"
	scheduleapi "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func main() {
	log, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer log.Sugar()

	bot, err := tele.NewBot(tele.Settings{
		Token:       "1604886434:AAHWVLB3R71nPz4VbVzvno3G8hKpTyWfylU",
		Poller:      &tele.LongPoller{Timeout: 2 * time.Second},
		Synchronous: false,
		Verbose:     true,
	})
	if err != nil {
		log.Fatal("cannot init tele.Bot", zap.Error(err))
	}

	groupsService, err := groups.NewFSFromFile("./data/groups.json")
	if err != nil {
		log.Fatal("cannot read groups file", zap.Error(err))
	}
	pool, err := pg.New(context.TODO(), pg.ConnString(
		"uksivt",
		"uksivt",
		"uksivt-schedule-bot",
		"",
		0,
	))
	if err != nil {
		log.Fatal("cannot connect to pg pool", zap.Error(err))
	}

	chatStorage := postgres.NewRepository(pool)
	chatService := chat.NewService(chatStorage)
	uksivtSchedule := scheduleapi.NewClient(nil, "https://back.uksivt.com/api/v1")
	handler := telegram.NewHandler(chatService, uksivtSchedule, groupsService, log)

	storage := memory.NewStorage()
	m := fsm.NewManager(bot.Group(), storage)
	handler.Route(m)
	bot.Start()
}
