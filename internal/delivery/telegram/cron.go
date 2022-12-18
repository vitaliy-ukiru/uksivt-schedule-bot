package telegram

import (
	"context"
	"fmt"
	"time"

	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/chat"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/scheduler"
	scheduleapi "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func (h Handler) CronJobSchedule() {
	now := time.Now()
	ctx := context.TODO()
	crons, err := h.crons.At(ctx, now)
	if err != nil {
		h.logger.With(zap.Error(err)).Error("cannot fetch crons")
	}

	for _, cron := range crons {
		c, err := h.uc.ByID(ctx, cron.ChatID)
		if err != nil {
			h.logger.Error("cannot get chat", zap.Error(err))
		}

		day := now
		if cron.Flags.Has(scheduler.NextDay) {
			forNextDay := 24 * time.Hour

			if day.Weekday() == time.Saturday {
				forNextDay *= 2
			}

			day = day.Add(forNextDay)
		}

		lessons, err := h.uksivt.LessonsOneDay(ctx, *c.Group, day)
		if err != nil {
			h.logger.Error("cannot fetch lessons", zap.Error(err))
			continue
		}

		params := &cronParams{
			Day:     day,
			Lessons: lessons,
			Cron:    cron,
			Chat:    c,
		}

		switch {
		case cron.Flags.Has(scheduler.FullOnlyIfReplaces):
			go h.cronFullOnReplace(params)

		case cron.Flags.HasAny(scheduler.ReplacesAlways, scheduler.OnlyIfHaveReplaces):
			go h.cronReplaces(params)

		case cron.Flags.Has(scheduler.Full):
			go h.bot.Send(
				tele.ChatID(c.TgID),
				lessonsToString(day, lessons),
			)
		}

		h.logger.With(zap.Int64("cron_id", cron.ID)).Debug("cron handled")

		//time.Sleep(300 * time.Millisecond)
	}

}

type cronParams struct {
	Day     time.Time
	Lessons []scheduleapi.Lesson
	Cron    scheduler.CronJob
	Chat    *chat.Chat
}

func (h Handler) cronFullOnReplace(p *cronParams) error {
	{
		var hasRepl bool
		for _, lesson := range p.Lessons {
			if lesson.Replacement {
				hasRepl = true
				break
			}
		}

		if !hasRepl {
			return nil
		}
	}

	_, err := h.bot.Send(tele.ChatID(p.Chat.TgID), lessonsToString(p.Day, p.Lessons))
	return err
}

func (h Handler) cronReplaces(p *cronParams) error {
	var replaces []scheduleapi.Lesson
	for _, lesson := range p.Lessons {
		if lesson.Replacement {
			replaces = append(replaces, lesson)
		}
	}

	if len(replaces) == 0 {
		if p.Cron.Flags.Has(scheduler.ReplacesAlways) {
			_, err := h.bot.Send(
				tele.ChatID(p.Chat.TgID),
				fmt.Sprintf("Замен на %s не найдено", p.Day.Format("02.01.2006")),
			)
			return err
		}
		return nil
	}

	_, err := h.bot.Send(
		tele.ChatID(p.Chat.TgID),
		lessonsToString(p.Day, replaces),
	)
	return err
}