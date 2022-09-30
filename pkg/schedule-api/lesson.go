package scheduleapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
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

func (l *LessonTimePair) UnmarshalJSON(bytes []byte) error {
	var data string
	if err := json.Unmarshal(bytes, &data); err != nil {
		return err
	}

	parts := strings.SplitN(data, " ", 2)
	start, err := time.Parse(lessonStartTimeFormat, parts[0])
	if err != nil {
		return err
	}

	end, err := time.Parse(lessonEndTimeFormat, parts[1])
	if err != nil {
		return err
	}

	//if start.After(end) {
	//	return errors.New("start time after end")
	//}

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

func (l *LessonTime) UnmarshalJSON(data []byte) error {
	if data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("invalid type for LessonTime")
	}

	msg := data[1 : len(data)-1]
	split := bytes.Split(msg, []byte{' '})
	if len(split)%2 != 0 {
		return errors.New("invalid count of time pairs")
	}
	buff := bytes.NewBuffer(nil)
	for i := 0; i < len(split); i += 2 {
		buff.Reset()

		start, end := append([]byte(nil), split[i]...), split[i+1]
		var pair LessonTimePair

		buff.WriteByte('"')
		buff.Write(start)
		buff.WriteByte(' ')
		buff.Write(end)
		buff.WriteByte('"')

		if err := json.NewDecoder(buff).Decode(&pair); err != nil {
			return err
		}
		*l = append(*l, pair)
	}

	return nil
}

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
	return fmt.Sprintf(
		"%d. %s - %s (%s)\n  %s",
		l.LessonNumber,
		l.Name,
		l.Teacher,
		l.LessonHall,
		l.Time.StringJoin("\n  "),
	)
}

type WeekOfLessons [7][]Lesson

var ErrInvalidDayNumber = errors.New("day number out of time.Weekday")

func SetToWeek(week map[string][]Lesson) (result WeekOfLessons, err error) {
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
