package group

import (
	"context"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type PostgresService struct {
	c PgConnection
}

type PgConnection interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
}

func NewPostgresService(c PgConnection) *PostgresService {
	return &PostgresService{c: c}
}

func (p *PostgresService) Years(ctx context.Context) ([]int, error) {
	const sql = `SELECT DISTINCT year from public.groups ORDER BY year`
	rows, err := p.c.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []int
	if err := sliceScan(rows, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (p *PostgresService) Specs(ctx context.Context, year int) ([]string, error) {
	const sql = `SELECT DISTINCT spec from public.groups WHERE year=$1 ORDER BY spec`
	rows, err := p.c.Query(ctx, sql, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []string
	if err := sliceScan(rows, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (p *PostgresService) Numbers(ctx context.Context, year int, spec string) ([]int, error) {
	const sql = `SELECT num from public.groups 
           		WHERE year=$1 AND spec=$2
           		ORDER BY spec`

	rows, err := p.c.Query(ctx, sql, year, spec)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []int
	if err := sliceScan(rows, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func sliceScan[T any](rows pgx.Rows, dest *[]T) error {
	var item T
	for rows.Next() {
		if err := rows.Scan(&item); err != nil {
			return err
		}
		*dest = append(*dest, item)
	}
	return nil
}
