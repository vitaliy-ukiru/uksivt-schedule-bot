package group

import (
	"context"

	scheduleapi "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
)

type Storage interface {
	FindByID(ctx context.Context, id int) (Group, error)
	FindID(ctx context.Context, params Group) (int, error)

	Years(ctx context.Context) ([]int, error)
	Specs(ctx context.Context, year int) ([]string, error)
	Nums(ctx context.Context, year int, spec string) ([]int, error)
}

type Usecase interface {
	FindID(ctx context.Context, group string) (int, error)

	ByID(ctx context.Context, id int) (string, error)
}

type Service struct {
	store Storage
}

func NewService(store Storage) *Service {
	return &Service{store: store}
}

func (s *Service) Years(ctx context.Context) ([]int, error) {
	return s.store.Years(ctx)
}

func (s *Service) Specs(ctx context.Context, year int) ([]string, error) {
	return s.store.Specs(ctx, year)
}

func (s *Service) Numbers(ctx context.Context, year int, spec string) ([]int, error) {
	return s.store.Nums(ctx, year, spec)
}

func (s *Service) FindID(ctx context.Context, group string) (int, error) {
	g, err := scheduleapi.ParseGroup(group)
	if err != nil {
		return 0, err
	}

	return s.store.FindID(ctx, Group(g))
}
func (s *Service) ByID(ctx context.Context, id int) (string, error) {
	group, err := s.store.FindByID(ctx, id)
	if err != nil {
		return "", err
	}
	return group.String(), err
}
