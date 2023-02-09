package postgres

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgtype"
	"github.com/mitchellh/mapstructure"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/scheduler"
)

type Repository struct {
	q      *DBQuerier
	conn   Connection
	period time.Duration
}

type Connection genericConn

func NewRepository(conn Connection) *Repository {
	return &Repository{q: NewQuerier(conn), conn: conn}
}

func (r *Repository) Insert(ctx context.Context, cron scheduler.CronJob) (int64, error) {
	return r.q.CreateJob(ctx, CreateJobParams{
		ChatID: cron.ChatID,
		SendAt: cron.At,
		Flags:  uint16(cron.Flags),
		Title: pgtype.Varchar{
			String: cron.Title,
			Status: pgtype.Present,
		},
	})
}

func (r *Repository) FindInPeriod(ctx context.Context, at time.Time, periodRange time.Duration) ([]scheduler.CronJob, error) {
	at = at.Round(time.Minute)

	rows, err := r.q.FindInPeriod(ctx, at, periodRange)
	if err != nil {
		return nil, err
	}

	result := make([]scheduler.CronJob, len(rows))
	for i, row := range rows {
		result[i] = cronRow(row).ToDomain()
	}
	return result, nil
}

func (r *Repository) FindAtTime(ctx context.Context, at time.Time) ([]scheduler.CronJob, error) {
	at = at.Round(time.Minute)

	rows, err := r.q.FindAtTime(ctx, at)
	if err != nil {
		return nil, err
	}

	result := make([]scheduler.CronJob, len(rows))
	for i, row := range rows {
		result[i] = cronRow(row).ToDomain()
	}
	return result, nil
}

func (r *Repository) FindByChat(ctx context.Context, chatId int64) ([]scheduler.CronJob, error) {
	rows, err := r.q.FindByChat(ctx, chatId)
	if err != nil {
		return nil, err
	}

	result := make([]scheduler.CronJob, len(rows))
	for i, row := range rows {
		result[i] = cronRow(row).ToDomain()
	}
	return result, nil
}

func (r *Repository) FindByID(ctx context.Context, chatId int64) (*scheduler.CronJob, error) {
	row, err := r.q.FindByID(ctx, chatId)
	if err != nil {
		return nil, err
	}
	cron := cronRow(row).ToDomain()
	return &cron, nil
}

const table = "crons"

func (r *Repository) Update(ctx context.Context, job scheduler.CronJob) error {
	var fields map[string]any

	if err := mapstructure.Decode(&job, &fields); err != nil {
		return err
	}
	if job.At.IsZero() {
		delete(fields, "send_at")
	} else {
		fields["send_at"] = job.At
	}
	sql, args, err := sq.
		Update(table).
		SetMap(fields).
		Where(sq.Eq{"id": job.ID}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.conn.Exec(ctx, sql, args...)
	return err
}

func (r *Repository) Delete(ctx context.Context, cronId int64) error {
	_, err := r.q.Delete(ctx, cronId)
	return err
}

type cronRow struct {
	ID     int64
	ChatID int64
	Title  pgtype.Varchar
	SendAt time.Time
	Flags  uint16
}

func (row cronRow) ToDomain() scheduler.CronJob {
	return scheduler.CronJob{
		ID:     row.ID,
		At:     row.SendAt,
		Title:  row.Title.String,
		Flags:  scheduler.FlagSet(row.Flags),
		ChatID: row.ChatID,
	}
}
