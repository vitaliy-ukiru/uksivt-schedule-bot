package chat

import (
	"time"

	scheduleapi "github.com/vitaliy-ukiru/uksivt-schedule-bot/pkg/schedule-api"
)

type Chat struct {
	ID        int64              `json:"id"`
	ChatID    int64              `json:"chat_id"`
	Group     *scheduleapi.Group `json:"group,omitempty"`
	CreatedAt time.Time          `json:"created_a,omitempty"`
	DeletedAt *time.Time         `json:"deleted_at,omitempty"`
}
