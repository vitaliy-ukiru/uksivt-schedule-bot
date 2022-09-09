package chat

import (
	"context"

	scheduleapi "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/uksivt-schedule-api"
)

type Usecase interface {
	Create(ctx context.Context, tgId int64) (*Chat, error)

	Lookup(ctx context.Context, tgId int64) (*Chat, error)
	ByTelegramID(ctx context.Context, chatTgId int64) (*Chat, error)
	SetGroup(ctx context.Context, chatTgID int64, group scheduleapi.Group) error
	ClearGroup(ctx context.Context, chatTgID int64) error
	Restore(ctx context.Context, chatTgID int64) (*Chat, error)
	Delete(ctx context.Context, chatTgID int64) error
}
