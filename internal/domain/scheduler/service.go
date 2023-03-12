package scheduler

import (
	"context"
	"time"
)

type Usecase interface {
	Create(ctx context.Context, dto CreateJobDTO) (*CronJob, error)

	ByID(ctx context.Context, cronId int64) (*CronJob, error)
	At(ctx context.Context, t time.Time) ([]CronJob, error)
	ForChat(ctx context.Context, chatId int64) ([]CronJob, error)
	CountInChat(ctx context.Context, chatId int64) (int64, error)

	Update(ctx context.Context, cron CronJob) error
	Delete(ctx context.Context, id int64) error
}

type Storage interface {
	Insert(ctx context.Context, base CronJob) (int64, error)

	FindInPeriod(ctx context.Context, at time.Time, periodRange time.Duration) ([]CronJob, error)
	FindAtTime(ctx context.Context, at time.Time) ([]CronJob, error)
	FindByChat(ctx context.Context, chatId int64) ([]CronJob, error)
	CountInChat(ctx context.Context, chatId int64) (int64, error)
	FindByID(ctx context.Context, cronId int64) (*CronJob, error)

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

func (s *Service) Create(ctx context.Context, dto CreateJobDTO) (*CronJob, error) {
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

func (s *Service) ByID(ctx context.Context, cronId int64) (*CronJob, error) {
	return s.store.FindByID(ctx, cronId)
}

func (s *Service) At(ctx context.Context, t time.Time) ([]CronJob, error) {
	crons, err := s.store.FindAtTime(ctx, t)
	if err != nil {
		return nil, err
	}

	return crons, nil
}

func (s *Service) ForChat(ctx context.Context, chatId int64) ([]CronJob, error) {
	crons, err := s.store.FindByChat(ctx, chatId)
	return crons, err
}

func (s *Service) InPeriod(ctx context.Context, t time.Time) ([]CronJob, error) {
	crons, err := s.store.FindInPeriod(ctx, t, s.period)
	if err != nil {
		return nil, err
	}

	return crons, nil
}

func (s *Service) CountInChat(ctx context.Context, chatId int64) (int64, error) {
	return s.store.CountInChat(ctx, chatId)
}

func (s *Service) Update(ctx context.Context, cron CronJob) error {
	return s.store.Update(ctx, cron)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}
