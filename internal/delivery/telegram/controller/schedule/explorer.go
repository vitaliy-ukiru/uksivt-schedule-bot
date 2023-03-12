package schedule

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ivahaev/russian-time"
	"github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/adapters/schedule"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/telegram/callback"

	tele "gopkg.in/telebot.v3"
)

func (h *Handler) LessonsCommand(c tele.Context, _ fsm.Context) error {
	chat := h.getChat(c.Chat().ID)
	if chat == nil {
		return c.Send("error: cannot get chat")
	}

	payload := c.Data()
	if chat.Group == nil && payload == "" {
		return c.Send(
			"Укажите группу через пробел после команды (Пример: /lessons 20П-1)" +
				"или установите её через /select_group",
		)
	}

	var group string
	{
		if chat.Group != nil {
			group = *chat.Group
		}

		if payload != "" {
			group = c.Data()

		}
	}

	t := time.Now()
	if t.Weekday() == time.Sunday {
		t = t.AddDate(0, 0, 1)
	}

	lessons, err := h.uksivt.LessonsOneDay(context.TODO(), group, t)
	if errors.Is(err, schedule.ErrInvalidGroup) {
		return c.Send(
			"Не получилось достать группу из аргумента команды.\n" +
				"Пример корректного ввода: 20П-1",
		)
	}
	if err != nil {
		return c.Send("error: " + err.Error())
	}
	return c.Send(lessonsToString(t, lessons), ExplorerMarkup(t, group))
}

func (h *Handler) ExplorerCallback(c tele.Context, data callback.M) error {
	day, err := time.Parse("2006-01-02", data["day"])
	if err != nil {
		return answerCallback(c, "invalid callback day", true)
	}

	group := data["g"]

	lessons, err := h.uksivt.LessonsOneDay(context.TODO(), group, day)
	if errors.Is(err, schedule.ErrInvalidGroup) {
		return answerCallback(c, "invalid callback group", true)
	}
	if err != nil {
		return answerCallback(c, "error: "+err.Error(), true)
	}

	return c.EditOrSend(lessonsToString(day, lessons), ExplorerMarkup(day, group))

}

func lessonsToString(day time.Time, lessons []schedule.Lesson) string {
	buff := make([]string, len(lessons)+1)

	lt := rtime.Time(day)
	buff[0] = fmt.Sprintf(
		"%d %s | %s",
		day.Day(),
		lt.Month().StringInCase(),
		lt.Weekday(),
	)
	for i, lesson := range lessons {
		buff[i+1] = lesson.StringReplacement("<b>[ЗАМЕНА]</b> ")
	}

	return strings.Join(buff, "\n\n")
}

func answerCallback(c tele.Context, text string, alert bool) error {
	return c.Respond(&tele.CallbackResponse{
		Text:      text,
		ShowAlert: alert,
	})
}
