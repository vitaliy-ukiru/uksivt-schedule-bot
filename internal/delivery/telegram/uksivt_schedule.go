package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	fsm "github.com/vitaliy-ukiru/fsm-telebot"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/delivery/telegram/keyboards"
	scheduleapi "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
	"github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/telegram/callback"
	tele "gopkg.in/telebot.v3"
)

func (h Handler) GetGroupCommand(c tele.Context, _ fsm.FSMContext) error {
	chat := getChat(c)
	if chat == nil {
		return c.Send("error: cannot get chat")
	}
	if chat.Group == nil {
		return c.Send("Группа для чата не установлена")
	}

	return c.Send(fmt.Sprintf("Для чата выбрана группа %s", chat.Group.String()))
}

func (h Handler) ScheduleCommand(c tele.Context, _ fsm.FSMContext) error {
	chat := getChat(c)
	if chat == nil {
		return c.Send("error: cannot get chat")
	}

	payload := c.Data()
	if chat.Group == nil && payload == "" {
		return c.Send("Укажите группу через пробел после команды " +
			"или установите её через /select_group")
	}

	var group scheduleapi.Group
	if chat.Group != nil {
		group = *chat.Group
	}
	if payload != "" {
		var err error
		group, err = scheduleapi.ParseGroup(c.Data())
		if err != nil {
			return c.Send("Не получилось достать группу из аргумента команды.\n" +
				"Пример корректного ввода: 22ОИБ-1")
		}
	}
	t := time.Now()
	lessons, err := h.getLessons(group, t)
	if err != nil {
		return c.Send("error: " + err.Error())
	}
	return c.Send(lessons, keyboards.ScheduleMarkup(t, group))
}

func (h Handler) ScheduleExplorerCallback(c tele.Context, data callback.M) error {
	day, err := time.Parse("2006-01-02", data["day"])
	if err != nil {
		return answerCallback(c, "invalid callback day", true)
	}
	group, err := scheduleapi.ParseGroup(data["g"])
	if err != nil {
		return answerCallback(c, "invalid callback group", true)
	}

	lessons, err := h.getLessons(group, day)
	if err != nil {
		return answerCallback(c, "error: "+err.Error(), true)
	}
	return c.EditOrSend(lessons, keyboards.ScheduleMarkup(day, group))

}

func (h Handler) getLessons(group scheduleapi.Group, today time.Time) (string, error) {
	setLessons, err := h.uksivt.Lessons(context.TODO(), group, getStartDayOfWeek(today))
	if err != nil {
		return "", err
	}

	lessons, err := scheduleapi.SetToWeek(setLessons)
	if err != nil {
		return "", err
	}

	lessonsToday := lessons[today.Weekday()-1]
	buff := make([]string, len(lessonsToday)+1)
	buff[0] = today.Format("January 02 | Monday")
	for i, lesson := range lessonsToday {
		buff[i+1] = lesson.String()
	}
	return strings.Join(buff, "\n\n"), nil
}

func answerCallback(c tele.Context, text string, alert bool) error {
	return c.Respond(&tele.CallbackResponse{
		Text:      text,
		ShowAlert: alert,
	})
}
