package scheduleapi

import (
	"errors"
	"fmt"
	"time"
)

type Lesson struct {
	Group        string     `json:"college_group"`
	DayOfWeek    int        `json:"day_of_week"`
	LessonNumber int        `json:"lesson_number"`
	Name         string     `json:"lesson"`
	Teacher      string     `json:"teacher"`
	LessonHall   string     `json:"lesson_hall"`
	Replacement  bool       `json:"replacement"`
	Time         LessonTime `json:"time"`
}

func (l Lesson) String() string {
	return l.StringReplacement("")
}

// StringReplacement returns string view of Lesson and add onReplacement
// before Lesson.Name. IT NOT ADD WHITESPACE. Just join.
func (l Lesson) StringReplacement(onReplacement string) string {
	lesson := l.Name
	if l.Replacement {
		lesson = onReplacement + lesson
	}
	return fmt.Sprintf(
		"%d. %s - %s (%s)\n  %s",
		l.LessonNumber,
		lesson,
		l.Teacher,
		l.LessonHall,
		l.Time.StringJoin("\n  "),
	)
}

type WeekOfLessons [6][]Lesson

var ErrInvalidDayNumber = errors.New("day number out of time.Weekday")

func SetToWeek(week map[time.Weekday][]Lesson) (WeekOfLessons, error) {
	var result WeekOfLessons
	for day, lessons := range week {
		if day < time.Monday || day > time.Saturday {
			return result, ErrInvalidDayNumber
		}

		// API and time.Weekday storages Monday as 1
		// but WeekOfLessons starts with 0 (like everything in CS)
		result[day-1] = lessons
	}
	return result, nil
}
