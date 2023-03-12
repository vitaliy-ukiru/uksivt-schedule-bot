package schedule

import (
	"fmt"
	"regexp"
	"strings"
)

type Lesson struct {
	LessonNumber int    `json:"lesson_number"`
	Name         string `json:"lesson"`
	Teacher      string `json:"teacher"`
	LessonHall   string `json:"lesson_hall"`
	Time         string `json:"time"`
	Replacement  bool   `json:"replacement"`
}

// StringReplacement returns string view of Lesson and add onReplacement
// before Lesson.Name. IT NOT ADD WHITESPACE. Just join.
func (l Lesson) StringReplacement(onReplacement string) string {
	var title strings.Builder
	if l.Replacement {
		title.Grow(len(onReplacement) + len(l.Name))
	}

	if l.Replacement {
		title.WriteString(onReplacement)
	}
	title.WriteString(l.Name)

	if l.Teacher != "" {
		title.WriteString(" - ")
		title.WriteString(l.Teacher)
	}

	if l.LessonHall != "" {
		title.WriteByte('(')
		title.WriteString(l.LessonHall)
		title.WriteByte(')')
	}

	return fmt.Sprintf(
		"%d. %s\n%s",
		l.LessonNumber,
		title.String(),
		l.fmtTime(),
	)
}

var reLessonTimes = regexp.MustCompile(`s(\d+:\d+)\se(\d+:\d+)`)

func (l Lesson) fmtTime() string {
	matches := reLessonTimes.FindAllStringSubmatch(l.Time, -1)
	formatted := make([]string, len(matches))

	for i, match := range matches {
		formatted[i] = fmt.Sprintf(" %s  — %s", match[1], match[2])
	}

	return strings.Join(formatted, "\n")
}
