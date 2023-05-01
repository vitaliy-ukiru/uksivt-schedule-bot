package postgres

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	. "github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat"
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

func (r *Repository) Create(ctx context.Context, chatId int64) (CreateChatDTO, error) {
	chatRow, err := r.q.CreateChat(ctx, chatId)
	if err != nil {
		return CreateChatDTO{}, errors.Wrap(err, "pg.create")
	}

	return CreateChatDTO{
		ID:        chatRow.ID,
		CreatedAt: chatRow.CreatedAt.Time,
	}, nil
}

func (r *Repository) FindByID(ctx context.Context, chatId int64) (*ModelDTO, error) {
	row, err := r.q.FindByID(ctx, chatId)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrChatNotFound
	}
	if err != nil {
		return nil, errors.Wrap(err, "pg.by_id")
	}
	return rowType(row).ToModel(), nil
}

func (r *Repository) FindByTelegramID(ctx context.Context, id int64) (*ModelDTO, error) {
	row, err := r.q.FindByTgID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrChatNotFound
		}
		return nil, errors.Wrap(err, "pg.by_tg")
	}
	return rowType(row).ToModel(), nil
}

func (r *Repository) FindAllActiveChats(ctx context.Context) ([]*ModelDTO, error) {
	rows, err := r.q.FindAllActiveChats(ctx)
	if err != nil {
		return nil, err
	}

	models := make([]*ModelDTO, len(rows))
	for i, row := range rows {
		models[i] = rowType(row).ToModel()
	}
	return models, nil
}

func (r *Repository) UpdateChatGroup(ctx context.Context, id int64, group *int) error {
	const updateGroupSQL = `UPDATE chats
SET group_id = $1
WHERE chat_id = $2;`

	tag, err := r.c.Exec(ctx, updateGroupSQL, group, id)

	if err != nil {
		return errors.Wrap(err, "pg.update: exec query UpdateGroup")
	}
	return checkRowsAffected(tag)
}

func (r *Repository) RestoreFromDeleted(ctx context.Context, id int64) error {
	tag, err := r.q.UndeleteChat(ctx, id)
	if err != nil {
		return err
	}
	return checkRowsAffected(tag)
}

func (r *Repository) Delete(ctx context.Context, chatId int64) error {
	tag, err := r.q.Delete(ctx, chatId)
	if err != nil {
		return errors.Wrap(err, "pg.delete")
	}
	return checkRowsAffected(tag)
}

func checkRowsAffected(tag pgconn.CommandTag) (err error) {
	if tag.RowsAffected() == 0 {
		err = ErrNotModified
	}
	return
}

func (r *Repository) Session(ctx context.Context, fn func(session Storage) error) error {
	return r.c.BeginFunc(ctx, func(tx pgx.Tx) error {
		withTx, _ := r.q.WithTx(tx)
		return fn(&Repository{q: withTx, c: tx})
	})
}

type rowType FindByIDRow

func (row rowType) ToModel() *ModelDTO {
	chat := &ModelDTO{
		ID:        row.ID,
		TgID:      row.ChatID,
		CreatedAt: row.CreatedAt.Time,
	}
	if row.DeletedAt.Status == pgtype.Present {
		chat.DeletedAt = &row.DeletedAt.Time
	}

	if row.GroupID != 0 {
		chat.Group = &row.GroupID
	}

	return chat
}
