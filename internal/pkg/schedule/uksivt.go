package schedule

import (
	"context"
	"time"

	"github.com/pkg/errors"
	api "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
)

type Usecase interface {
	LessonsForWeek(ctx context.Context, group api.Group, weekStart time.Time) (api.WeekOfLessons, error)
	LessonsOneDay(ctx context.Context, group api.Group, today time.Time) ([]api.Lesson, error)
}

type Service struct {
	c *api.Client
}

func NewService(c *api.Client) *Service {
	return &Service{c: c}
}

func (s Service) LessonsForWeek(ctx context.Context, group api.Group, weekStart time.Time) (api.WeekOfLessons, error) {
	lessonsSet, err := s.c.Lessons(ctx, group, weekStart)
	if err != nil {
		return api.WeekOfLessons{}, errors.Wrap(err, "cannot fetch lessons")
	}
	return api.SetToWeek(lessonsSet)
}

func (s Service) LessonsOneDay(ctx context.Context, group api.Group, today time.Time) ([]api.Lesson, error) {
	wd := today.Weekday()
	if wd < time.Monday || wd > time.Saturday {
		return nil, errors.New("it not study day")
	}

	week, err := s.c.Lessons(ctx, group, getStartDayOfWeek(today))
	if err != nil {
		return nil, err
	}

	lessons := week[wd]
	return lessons, nil
}

func getStartDayOfWeek(tm time.Time) time.Time { //get monday 00:00:00
	weekday := time.Duration(tm.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	year, month, day := tm.Date()
	currentZeroDay := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
	return currentZeroDay.Add(-1 * (weekday - 1) * 24 * time.Hour)
}
