package chat

import (
	"context"

	"github.com/pkg/errors"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/group"
)

type LookupStatus byte

const (
	StatusNone     LookupStatus = iota // undefined or null
	StatusCreated                      // chat created now
	StatusFound                        // chat found in db
	StatusRestored                     // chat was deleted before lookup
)

type Usecase interface {
	Create(ctx context.Context, tgId int64) (*Chat, error)

	ByID(ctx context.Context, chatId int64) (*Chat, error)
	Lookup(ctx context.Context, tgId int64) (*Chat, LookupStatus, error)
	ByTelegramID(ctx context.Context, chatTgId int64) (*Chat, error)
	//ActiveChats(ctx context.Context) ([]Chat, error)

	SetGroup(ctx context.Context, chatTgID int64, group string) error
	ClearGroup(ctx context.Context, chatTgID int64) error

	Restore(ctx context.Context, chatTgID int64) (*Chat, error)
	Delete(ctx context.Context, chatTgID int64) error
}

type Storage interface {
	Create(ctx context.Context, chatId int64) (CreateChatDTO, error)
	FindByID(ctx context.Context, chatId int64) (*ModelDTO, error)
	FindByTelegramID(ctx context.Context, id int64) (*ModelDTO, error)

	UpdateChatGroup(ctx context.Context, id int64, group *int) error
	RestoreFromDeleted(ctx context.Context, id int64) error

	Delete(ctx context.Context, chatId int64) error
	Session(ctx context.Context, fn func(session Storage) error) error
}

type Service struct {
	store  Storage
	groups group.Usecase
}

func NewService(store Storage, groups group.Usecase) *Service {
	return &Service{store: store, groups: groups}
}

func (s *Service) Create(ctx context.Context, tgId int64) (*Chat, error) {
	dto, err := s.store.Create(ctx, tgId)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create chat id=%d", tgId)
	}
	return &Chat{
		ID:        dto.ID,
		TgID:      tgId,
		CreatedAt: dto.CreatedAt,
	}, nil

}

func (s *Service) Lookup(ctx context.Context, tgId int64) (*Chat, LookupStatus, error) {
	chat, err := s.ByTelegramID(ctx, tgId)
	if err == nil {
		return chat, StatusFound, nil
	}

	if errors.Is(err, ErrChatNotFound) {
		chat, err = s.Create(ctx, tgId)

		var status LookupStatus
		if err == nil {
			status = StatusCreated
		}
		return chat, status, err
	}

	if errors.Is(err, ErrChatDeleted) {
		// in this case Usecase returning *Chat and err ErrChatDeleted
		// see Service.ByTelegramID method.
		// For optimize db queries I make only restore query
		// i.e. I already have chat instance.
		if err := s.store.RestoreFromDeleted(ctx, tgId); err != nil {
			return nil, StatusNone, err
		}
		//TODO: add logs when chat restored
		chat.DeletedAt = nil

		return chat, StatusRestored, nil
	}
	return nil, StatusNone, err
}

func (s *Service) ByID(ctx context.Context, chatId int64) (*Chat, error) {
	model, err := s.store.FindByID(ctx, chatId)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find_by_id id=%d", chatId)
	}
	return s.modelToDomain(ctx, model)
}

func (s *Service) ByTelegramID(ctx context.Context, chatTgId int64) (*Chat, error) {
	model, err := s.store.FindByTelegramID(ctx, chatTgId)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find_by_tg chat=%d", chatTgId)
	}
	return s.modelToDomain(ctx, model)
}

func (s *Service) modelToDomain(ctx context.Context, model *ModelDTO) (*Chat, error) {
	chat := model.chat()
	if chat.IsDeleted() {
		return chat, ErrChatDeleted
	}
	if model.Group != nil {
		g, err := s.groups.ByID(ctx, *model.Group)
		if err != nil {
			return nil, err
		}
		chat.Group = &g
	}
	return chat, nil
}

func (s *Service) SetGroup(ctx context.Context, chatTgID int64, group string) error {
	groupId, err := s.groups.FindID(ctx, group)
	if err != nil {
		return err
	}

	return errors.Wrapf(
		s.store.UpdateChatGroup(ctx, chatTgID, &groupId),
		"cannot set group chat=%d group=%+v", chatTgID, group,
	)
}

func (s *Service) ClearGroup(ctx context.Context, chatTgID int64) error {
	return errors.Wrapf(
		s.store.UpdateChatGroup(ctx, chatTgID, nil),
		"cannot delete group chat=%d", chatTgID,
	)
}

func (s *Service) Restore(ctx context.Context, chatTgID int64) (*Chat, error) {
	var chat *Chat
	err := s.store.Session(ctx, func(store Storage) error {
		var err error
		if err = store.RestoreFromDeleted(ctx, chatTgID); err != nil {
			return errors.Wrapf(err, "cannot restore chat=%d", chatTgID)
		}

		chat, err = s.withSession(store).ByTelegramID(ctx, chatTgID)
		return errors.Wrapf(err, "cannot find restored chat=%d", chatTgID)
	})

	return chat, err

}
func (s *Service) withSession(session Storage) *Service {
	return &Service{store: session}
}

func (s *Service) Delete(ctx context.Context, chatTgID int64) error {
	return errors.Wrapf(
		s.store.Delete(ctx, chatTgID),
		"cannot delete chat=%d", chatTgID,
	)
}
