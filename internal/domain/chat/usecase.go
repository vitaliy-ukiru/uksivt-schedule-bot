package chat

import (
	"context"
	"errors"

	scheduleapi "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
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

type Storage interface {
	Create(ctx context.Context, chatId int64) (CreateChatDTO, error)
	FindByTelegramID(ctx context.Context, id int64) (*Chat, error)
	UpdateChatGroup(ctx context.Context, id int64, group *scheduleapi.Group) error
	RestoreFromDeleted(ctx context.Context, id int64) error
	Delete(ctx context.Context, chatId int64) error
}

type Service struct {
	store Storage
}

func NewService(store Storage) Usecase {
	return &Service{store: store}
}

func (s Service) Create(ctx context.Context, tgId int64) (*Chat, error) {
	dto, err := s.store.Create(ctx, tgId)
	if err != nil {
		return nil, err
	}
	return &Chat{
		ID:        dto.ID,
		ChatID:    tgId,
		CreatedAt: dto.CreatedAt,
	}, nil

}

func (s Service) Lookup(ctx context.Context, tgId int64) (*Chat, error) {
	chat, err := s.ByTelegramID(ctx, tgId)
	if err != nil {
		if errors.Is(err, ErrChatNotFound) {
			return s.Create(ctx, tgId)
		}

		if errors.Is(err, ErrChatDeleted) {
			// in this case Usecase returning *Chat and err ErrChatDeleted
			// see Service.ByTelegramID method
			if err := s.store.RestoreFromDeleted(ctx, tgId); err != nil {
				return nil, err
			}
			//TODO: add logs whet chat restored
			chat.DeletedAt = nil
			return chat, nil
		}

		return nil, err
	}
	return chat, nil
}

func (s Service) ByTelegramID(ctx context.Context, chatTgId int64) (*Chat, error) {
	chat, err := s.store.FindByTelegramID(ctx, chatTgId)
	if err != nil {
		return nil, err
	}
	if chat.DeletedAt != nil {
		return chat, ErrChatDeleted
	}
	return chat, nil
}

func (s Service) SetGroup(ctx context.Context, chatTgID int64, group scheduleapi.Group) error {
	return s.store.UpdateChatGroup(ctx, chatTgID, &group)
}

func (s Service) ClearGroup(ctx context.Context, chatTgID int64) error {
	return s.store.UpdateChatGroup(ctx, chatTgID, nil)
}

func (s Service) Restore(ctx context.Context, chatTgID int64) (*Chat, error) {
	if err := s.store.RestoreFromDeleted(ctx, chatTgID); err != nil {
		return nil, err
	}

	return s.ByTelegramID(ctx, chatTgID)
}

func (s Service) Delete(ctx context.Context, chatTgID int64) error {
	return s.store.Delete(ctx, chatTgID)
}
