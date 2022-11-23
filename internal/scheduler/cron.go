package scheduler

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat"
)

type CronFetcher interface {
	At(ctx context.Context, t time.Time) ([]CronJob, error)
}

type ChatUC interface {
	ByID(ctx context.Context, chatId int64) (*chat.Chat, error)
}

type Storage interface {
	Insert(ctx context.Context, base CronJobBase) (int64, error)
	FindForTime(ctx context.Context, at time.Time, periodRange time.Duration) ([]CronJobBase, error)
	Delete(ctx context.Context, cronId int64) error
}

type Service struct {
	store  Storage
	period time.Duration
	chatUC ChatUC
}

func (s Service) At(ctx context.Context, t time.Time) ([]CronJob, error) {
	jobBases, err := s.store.FindForTime(ctx, t, s.period)
	if err != nil {
		return nil, err
	}

	//chats, err := s.chatUC.ByIDMany(ctx, ids)
	if err != nil {
		return nil, err
	}

	crons := make([]CronJob, len(jobBases))
	for i, job := range jobBases {
		c, err := s.chatUC.ByID(ctx, job.ID)
		if err != nil {
			//TODO: delete crons if chat not accessibility
			return nil, errors.Wrapf(err, "chat_id=%d", job.ChatID)
		}

		crons[i] = CronJob{
			ID:    job.ID,
			At:    job.At,
			Flags: job.Flags,
			Chat:  c,
		}
	}

	return crons, nil

}
