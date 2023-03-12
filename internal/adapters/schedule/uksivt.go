package schedule

import (
	"context"
	"time"

	"github.com/pkg/errors"
	api "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
)

type Usecase interface {
	LessonsForWeek(ctx context.Context, group string, weekStart time.Time) (api.WeekOfLessons, error)
	LessonsOneDay(ctx context.Context, group string, today time.Time) ([]api.Lesson, error)
}

type Service struct {
	c *api.Client
}

func NewService(c *api.Client) *Service {
	return &Service{c: c}
}

var ErrInvalidGroup = api.ErrInvalidGroup

func (s Service) LessonsForWeek(ctx context.Context, group string, weekStart time.Time) (api.WeekOfLessons, error) {
	g, err := api.ParseGroup(group)
	if err != nil {
		return api.WeekOfLessons{}, ErrInvalidGroup
	}
	lessonsSet, err := s.c.Lessons(ctx, g, weekStart)

	if err != nil {
		return api.WeekOfLessons{}, errors.Wrap(err, "cannot fetch lessons")
	}
	return api.SetToWeek(lessonsSet)
}

func (s Service) LessonsOneDay(ctx context.Context, group string, today time.Time) ([]api.Lesson, error) {
	g, err := api.ParseGroup(group)
	if err != nil {
		return nil, errors.Wrap(err, "invalid group")
	}

	wd := today.Weekday()
	if wd < time.Monday || wd > time.Saturday {
		return nil, errors.New("it not study day")
	}

	week, err := s.c.Lessons(ctx, g, today)
	if err != nil {
		return nil, err
	}

	lessons := week[wd]
	return lessons, nil
}
