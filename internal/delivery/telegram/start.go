package telegram

import (
	"context"
	"fmt"
	"time"

	pkg "github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) StartCommand(c tele.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	chat, status, err := h.uc.Lookup(ctx, c.Chat().ID)
	if err != nil {
		h.logger.Error("cannot get chat", zap.Error(err))
		return c.Send("error: cannot get chat")
	}

	switch status {
	case pkg.StatusRestored:
		return c.Send(fmt.Sprintf(
			"Чат был удалён, возможно меня заблокировали. Но я восстановил данные.\n"+
				"В чате выбрана группа <i>%s<i>. Отправьте /select_group для изменения.", *chat.Group,
		))
	case pkg.StatusCreated:
		return c.Send("Я тут новенький, сохраняю чат в базу.\n" +
			"Отправьте /select_group для выбора группы.")
	case pkg.StatusFound:
		return c.Send("Зачем стартовать снова? Я уже и так тут есть")
	default:
		return c.Send("Возможно что-то пошло не по плану 🤨")
	}
}

func (h *Handler) HelpCommand(c tele.Context) error {
	const helpCmdText = `
Бот для отправки расписания для УКСИВТ. Имеет функциональность отправки расписания по указанному времени.
Работает в группах и личных чатах.
Информация основана на данных с https://uksivt.com
Разработчик: @ukirug
Исходный код: https://github.com/vitaliy-ukiru/uksivt-schedule-bot
Версия: 0.3-beta

<i>Команды:</i>
/select_group - Выбрать группу для чата.
/group - Посмотреть текущую группу в чате.
/lessons <code>[группа]</code> - Просмотреть расписание. Если после команды указать группы отправит для нее. По умолчанию используют группу, заданную через /select_group. 
/create - Создать программу отправки расписания (крон задача).
/crons - Список задач для текущего чата.
/edit - Панель изменения программ.`

	return c.Send(helpCmdText)
}
