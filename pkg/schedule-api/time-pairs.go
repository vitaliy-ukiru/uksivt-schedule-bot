package scheduleapi

import (
	"bytes"
	"encoding/json"
	"errors"
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

func (l *LessonTime) UnmarshalText(text []byte) error {
	split := bytes.Split(text, []byte{' '})
	if len(split)%2 != 0 {
		return errors.New("invalid count of time pairs")
	}
	buff := bytes.NewBuffer(nil)

	for i := 0; i < len(split); i += 2 {
		buff.Reset()

		start, end := split[i], split[i+1]

		{
			buff.WriteByte('"')
			buff.Write(start)

			buff.WriteByte(' ')

			buff.Write(end)
			buff.WriteByte('"')
		}

		var pair LessonTimePair
		if err := json.NewDecoder(buff).Decode(&pair); err != nil {
			return err
		}
		*l = append(*l, pair)
	}

	return nil
}

func unsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
