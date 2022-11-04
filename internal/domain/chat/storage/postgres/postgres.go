package postgres

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	domain "github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat"
	scheduleapi "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
)

type Repository struct {
	q *DBQuerier
	c Connection
}

type Connection interface {
	genericConn
	BeginFunc(ctx context.Context, f func(pgx.Tx) error) error
}

func NewRepository(conn Connection) *Repository {
	return &Repository{q: NewQuerier(conn), c: conn}
}

func (r Repository) Create(ctx context.Context, chatId int64) (domain.CreateChatDTO, error) {
	chatRow, err := r.q.CreateChat(ctx, chatId)
	if err != nil {
		return domain.CreateChatDTO{}, errors.Wrap(err, "pg.create")
	}

	return domain.CreateChatDTO{
		ID:        chatRow.ID,
		CreatedAt: chatRow.CreatedAt.Time,
	}, nil
}

func (r Repository) FindByTelegramID(ctx context.Context, id int64) (*domain.Chat, error) {
	row, err := r.q.FindByTgID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrChatNotFound
		}
		return nil, errors.Wrap(err, "pg.by_tg")
	}
	chat := &domain.Chat{
		ID:        row.ID,
		ChatID:    row.ChatID,
		CreatedAt: row.CreatedAt.Time,
	}
	if row.DeletedAt.Status == pgtype.Present {
		chat.DeletedAt = &row.DeletedAt.Time
	}

	if row.CollegeGroup.Status == pgtype.Present {
		g, err := scheduleapi.ParseGroup(row.CollegeGroup.String)
		if err != nil {
			return nil, errors.Wrap(err, "pg.by_tg.parse_group")
		}
		chat.Group = &g
	}

	return chat, nil
}

func (r Repository) UpdateChatGroup(ctx context.Context, id int64, group *scheduleapi.Group) error {
	g := pgtype.Text{
		Status: pgtype.Null,
	}
	if group != nil {
		g.String = group.String()
		g.Status = pgtype.Present
	}
	tag, err := r.q.UpdateGroup(ctx, g, id)
	if err != nil {
		return errors.Wrap(err, "pg.update")
	}
	return checkRowsAffected(tag)
}

func (r Repository) RestoreFromDeleted(ctx context.Context, id int64) error {
	tag, err := r.q.UndeleteChat(ctx, id)
	if err != nil {
		return err
	}
	return checkRowsAffected(tag)
}

func (r Repository) Delete(ctx context.Context, chatId int64) error {
	tag, err := r.q.Delete(ctx, chatId)
	if err != nil {
		return errors.Wrap(err, "pg.delete")
	}
	return checkRowsAffected(tag)
}

func checkRowsAffected(tag pgconn.CommandTag) (err error) {
	if tag.RowsAffected() == 0 {
		err = domain.ErrNotModified
	}
	return
}

func (r Repository) Session(ctx context.Context, fn func(session domain.Storage) error) error {
	return r.c.BeginFunc(ctx, func(tx pgx.Tx) error {
		withTx, _ := r.q.WithTx(tx)
		return fn(&Repository{q: withTx, c: tx})
	})
}
