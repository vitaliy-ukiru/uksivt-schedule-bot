package postgres

import (
	"context"

	"github.com/jackc/pgtype"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/group"
)

type Connection genericConn

type Repository struct {
	q *DBQuerier
	c Connection
}

func NewRepository(conn Connection) *Repository {
	return &Repository{q: NewQuerier(conn), c: conn}
}

func (r *Repository) FindByID(ctx context.Context, id int) (group.Group, error) {
	row, err := r.q.GroupByID(ctx, id)
	if err != nil {
		return group.Group{}, err
	}

	return group.Group{
		Year:   row.Year,
		Spec:   row.Spec.String,
		Number: row.Num,
	}, nil
}

func (r *Repository) FindID(ctx context.Context, params group.Group) (int, error) {
	groupId, err := r.q.IDByGroup(ctx, IDByGroupParams{
		Year: params.Year,
		Spec: pgtype.Text{String: params.Spec, Status: pgtype.Present},
		Num:  params.Number,
	})
	return groupId, err
}

func (r *Repository) Years(ctx context.Context) ([]int, error) {
	years, err := r.q.SelectYears(ctx)
	return years, err
}

func (r *Repository) Specs(ctx context.Context, year int) ([]string, error) {
	specs, err := r.q.SelectSpecsForYear(ctx, year)
	if err != nil {
		return nil, err
	}
	result := make([]string, len(specs))
	for i, spec := range specs {
		result[i] = spec.String
	}
	return result, nil
}

func (r *Repository) Nums(ctx context.Context, year int, spec string) ([]int, error) {
	nums, err := r.q.SelectNumsForYearAndSpec(
		ctx,
		year,
		pgtype.Text{String: spec, Status: pgtype.Present},
	)
	return nums, err
}
