package postgres

import (
	"context"

	"github.com/jackc/pgtype"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/group"
)

type Connection genericConn

type Repository struct {
	q *DBQuerier
	c Connection
}

func NewRepository(conn Connection) *Repository {
	return &Repository{q: NewQuerier(conn), c: conn}
}

func pgInt(x int) pgtype.Int2 {
	return pgtype.Int2{Int: int16(x), Status: pgtype.Present}
}

func (r *Repository) FindByID(ctx context.Context, id int16) (group.Group, error) {
	row, err := r.q.GroupByID(ctx, pgtype.Int2{Int: id, Status: pgtype.Present})
	if err != nil {
		return group.Group{}, err
	}

	return group.Group{
		Year:   int(row.Year.Int),
		Spec:   row.Spec.String,
		Number: int(row.Num.Int),
	}, nil
}

func (r *Repository) FindID(ctx context.Context, params group.Group) (int, error) {
	groupId, err := r.q.IDByGroup(ctx, IDByGroupParams{
		Year: pgInt(params.Year),
		Spec: pgtype.Text{String: params.Spec, Status: pgtype.Present},
		Num:  pgInt(params.Number),
	})
	return groupId, err
}

func (r *Repository) Years(ctx context.Context) ([]int, error) {
	years, err := r.q.SelectYears(ctx)
	return years, err
}

func (r *Repository) Specs(ctx context.Context, year int) ([]string, error) {
	specs, err := r.q.SelectSpecsForYear(ctx, pgInt(year))
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
		pgInt(year),
		pgtype.Text{String: spec, Status: pgtype.Present},
	)
	return nums, err
}
