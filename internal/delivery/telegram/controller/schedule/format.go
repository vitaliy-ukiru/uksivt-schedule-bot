package schedule

import (
	"fmt"
	"strings"
	"time"

	rtime "github.com/ivahaev/russian-time"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/adapters/schedule"
)

func lessonsToString(day time.Time, group string, lessons []schedule.Lesson) string {
	buff := make([]string, len(lessons)+1)

	lt := rtime.Time(day)
	buff[0] = fmt.Sprintf(
		"%d %s | %s | %s",
		day.Day(),
		lt.Month().StringInCase(),
		lt.Weekday(),
		group,
	)
	for i, lesson := range lessons {
		buff[i+1] = lesson.StringReplacement("<b>[ЗАМЕНА]</b> ")
	}

	return strings.Join(buff, "\n\n")
}
