package group

import (
	"context"
)

type Fetcher interface {
	Years(ctx context.Context) ([]int, error)
	Specs(ctx context.Context, year int) ([]string, error)
	Numbers(ctx context.Context, year int, spec string) ([]int, error)
}
