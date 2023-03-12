package chat

import (
	"time"
)

// TODO PROPOSAL[CHAT_PK]: delete field ID, use telegram id as ID

type Chat struct {
	ID        int64      `json:"id"`
	TgID      int64      `json:"tg_id"`
	Group     *string    `json:"group,omitempty"`
	CreatedAt time.Time  `json:"created_a,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

func (c Chat) IsDeleted() bool {
	return c.DeletedAt != nil
}
