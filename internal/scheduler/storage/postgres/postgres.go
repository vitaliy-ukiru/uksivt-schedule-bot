package postgres

import (
	"context"
	"time"

	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/scheduler"
)

type Repository struct {
	q      *DBQuerier
	period time.Duration
}

type Connection genericConn

func NewRepository(conn Connection) *Repository {
	return &Repository{q: NewQuerier(conn)}
}

func (r Repository) Insert(ctx context.Context, cron scheduler.CronJobBase) (int64, error) {
	return r.q.CreateJob(ctx, CreateJobParams{
		ChatID: cron.ChatID,
		SendAt: cron.At,
		Flags:  int16(cron.Flags),
	})
}

func (r Repository) FindForTime(ctx context.Context, at time.Time, periodRange time.Duration) ([]scheduler.CronJobBase, error) {
	at = at.Round(time.Minute).In(time.UTC)

	rows, err := r.q.FindAtTime(ctx, at, periodRange)
	if err != nil {
		return nil, err
	}

	result := make([]scheduler.CronJobBase, len(rows))
	for i, row := range rows {
		result[i] = scheduler.CronJobBase{
			ID:     row.ID,
			At:     row.SendAt,
			Flags:  scheduler.FlagSet(*row.Flags),
			ChatID: row.ChatID,
		}
	}
	return result, nil
}

func (r Repository) Delete(ctx context.Context, cronId int64) error {
	_, err := r.q.Delete(ctx, cronId)
	return err
}
