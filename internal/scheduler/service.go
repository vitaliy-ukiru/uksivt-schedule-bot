package scheduler

import (
	"context"
	"time"
)

type CronFetcher interface {
	At(ctx context.Context, t time.Time) ([]CronJob, error)
}

type Usecase interface {
	CronFetcher
	Create(ctx context.Context, dto CreateJobDTO) (*CronJob, error)
	Delete(ctx context.Context, id int64) error
}

type Storage interface {
	Insert(ctx context.Context, base CronJob) (int64, error)

	FindInPeriod(ctx context.Context, at time.Time, periodRange time.Duration) ([]CronJob, error)
	FindAtTime(ctx context.Context, at time.Time) ([]CronJob, error)
	FindByChat(ctx context.Context, chatId int64) ([]CronJob, error)

	Update(ctx context.Context, job CronJob) error
	Delete(ctx context.Context, cronId int64) error
}

type Service struct {
	store  Storage
	period time.Duration
}

func NewService(store Storage, period time.Duration) *Service {
	return &Service{store: store, period: period}
}

func (s Service) Create(ctx context.Context, dto CreateJobDTO) (*CronJob, error) {
	job := CronJob{
		Title:  dto.Title,
		At:     dto.At,
		Flags:  FlagSet(dto.Flags),
		ChatID: dto.ChatID,
	}
	var err error
	job.ID, err = s.store.Insert(ctx, job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (s Service) At(ctx context.Context, t time.Time) ([]CronJob, error) {
	crons, err := s.store.FindAtTime(ctx, t)
	if err != nil {
		return nil, err
	}

	return crons, nil
}

func (s Service) InPeriod(ctx context.Context, t time.Time) ([]CronJob, error) {
	crons, err := s.store.FindInPeriod(ctx, t, s.period)
	if err != nil {
		return nil, err
	}

	return crons, nil
}

func (s Service) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}
