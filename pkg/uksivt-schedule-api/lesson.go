package scheduleapi

import (
	"errors"
	"fmt"
	"strconv"
	"time"
)

type Lesson struct {
	Group        string `json:"college_group"`
	DayOfWeek    int    `json:"day_of_week"`
	LessonNumber int    `json:"lesson_number"`
	Name         string `json:"lesson"`
	Teacher      string `json:"teacher"`
	LessonHall   string `json:"lesson_hall"`
	Replacement  bool   `json:"replacement"`
	Time         string `json:"time"`
}

func (l Lesson) String() string {
	return fmt.Sprintf(
		"%s: %s - %s(%s)",
		l.Group,
		l.Teacher,
		l.Name,
		l.LessonHall,
	)
}

type WeekOfLessons [7][]Lesson

var ErrInvalidDayNumber = errors.New("day number out of time.Weekday")

func setToWeek(week map[string][]Lesson) (result WeekOfLessons, err error) {
	var dayInt int
	for day, lessons := range week {
		dayInt, err = strconv.Atoi(day)
		if err != nil {
			return
		}
		weekDay := time.Weekday(dayInt)
		dayInt-- // in api days starts from 1
		if weekDay < time.Monday || weekDay > time.Saturday {
			err = ErrInvalidDayNumber
			return
		}

		result[dayInt] = lessons
	}
	return
}
