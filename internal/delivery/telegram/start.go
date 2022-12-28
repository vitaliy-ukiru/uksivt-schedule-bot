package telegram

import (
	"context"
	"fmt"
	"time"

	pkg "github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/chat"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
)

func (h Handler) StartCommand(c tele.Context) error {
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
				"В чате выбрана группа <i>%s<i>. Отправьте /select_group для изменения.", chat.Group,
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
