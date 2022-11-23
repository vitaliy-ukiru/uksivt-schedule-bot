package telegram

import (
	"context"
	"fmt"
	"time"

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
		day := now
		if cron.Flags.Has(scheduler.NextDay) {
			forNextDay := 24 * time.Hour

			if day.Weekday() == time.Saturday {
				forNextDay *= 2
			}

			day = day.Add(forNextDay)
		}

		lessons, err := h.uksivt.LessonsOneDay(ctx, *cron.Chat.Group, day)
		if err != nil {
			h.logger.Error("cannot fetch lessons", zap.Error(err))
			continue
		}

		chatId := tele.ChatID(cron.Chat.ID)

		switch {
		case cron.Flags.Has(scheduler.Full):
			go h.bot.Send(
				chatId,
				lessonsToString(day, lessons),
			)
		case cron.Flags.HasAny(scheduler.ReplacesAlways, scheduler.OnlyIfHaveReplaces):
			go h.cronReplaces(day, lessons, cron)

		case cron.Flags.Has(scheduler.FullOnlyIfReplaces):
			go h.cronFullOnReplace(day, lessons, cron)
		}

		h.logger.With(zap.Int64("cron_id", cron.ID)).Debug("cron handled")

		//time.Sleep(300 * time.Millisecond)
	}

}

func (h Handler) cronFullOnReplace(day time.Time, lessons []scheduleapi.Lesson, cron scheduler.CronJob) error {
	{
		var hasRepl bool
		for _, lesson := range lessons {
			if lesson.Replacement {
				hasRepl = true
				break
			}
		}

		if !hasRepl {
			return nil
		}
	}

	_, err := h.bot.Send(tele.ChatID(cron.Chat.ChatID), lessonsToString(day, lessons))
	return err
}

func (h Handler) cronReplaces(day time.Time, lessons []scheduleapi.Lesson, cron scheduler.CronJob) error {
	var replaces []scheduleapi.Lesson
	for _, lesson := range lessons {
		if lesson.Replacement {
			replaces = append(replaces, lesson)
		}
	}

	if len(replaces) == 0 && cron.Flags.Has(scheduler.ReplacesAlways) {
		_, err := h.bot.Send(
			tele.ChatID(cron.Chat.ID),
			fmt.Sprintf("Замен на %s не найдено", day.Format("02.01.2006")),
		)
		return err
	}

	if len(replaces) == 0 {
		return nil
	}

	_, err := h.bot.Send(
		tele.ChatID(cron.Chat.ID),
		lessonsToString(day, lessons),
	)
	return err
}
