package scheduleapi

import (
	"strings"
	"time"
	"unsafe"
)

type LessonTimePair struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

const (
	lessonStartTimeFormat = "s15:04"
	lessonEndTimeFormat   = "e15:04"
)

func (l LessonTimePair) String() string {
	return l.Start.Format(lessonStartTimeFormat) + " " + l.End.Format(lessonEndTimeFormat)
}

func unsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func (l LessonTimePair) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *LessonTimePair) UnmarshalText(text []byte) error {
	parts := strings.SplitN(unsafeBytesToString(text), " ", 2)
	start, err := time.Parse(lessonStartTimeFormat, parts[0])
	if err != nil {
		return err
	}

	end, err := time.Parse(lessonEndTimeFormat, parts[1])
	if err != nil {
		return err
	}

	l.Start, l.End = start, end
	return nil
}

type LessonTime []LessonTimePair

func (l LessonTime) String() string {
	return l.StringJoin(" ")
}

func (l LessonTime) StringJoin(sep string) string {
	result := make([]string, len(l))
	for i, pair := range l {
		result[i] = pair.String()
	}
	return strings.Join(result, sep)
}
