package schedule

import (
	"context"
	"time"

	api "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
)

type Usecase interface {
	LessonsOneDay(ctx context.Context, group string, today time.Time) ([]Lesson, error)
}

type Service struct {
	c *api.Client
}

func NewService(c *api.Client) *Service {
	return &Service{c: c}
}

var ErrInvalidGroup = api.ErrInvalidGroup

func (s Service) LessonsOneDay(ctx context.Context, group string, today time.Time) ([]Lesson, error) {
	wd := today.Weekday()

	week, err := s.c.Lessons(ctx, group, today)
	if err != nil {
		return nil, err
	}

	lessons := week[wd]

	result := make([]Lesson, len(lessons))
	for i, l := range lessons {
		result[i] = Lesson{
			LessonNumber: l.LessonNumber,
			Name:         l.Name,
			Teacher:      l.Teacher,
			LessonHall:   l.LessonHall,
			Time:         l.Time,
			Replacement:  l.Replacement,
		}
	}
	return result, nil
}
