package chat

import "time"

type CreateChatDTO struct {
	ID        int64
	CreatedAt time.Time
}

type ModelDTO struct {
	ID        int64
	TgID      int64
	Group     *int16
	CreatedAt time.Time
	DeletedAt *time.Time
}

func (m ModelDTO) chat() *Chat {
	return &Chat{
		ID:        m.ID,
		TgID:      m.TgID,
		CreatedAt: m.CreatedAt,
		DeletedAt: m.DeletedAt,
	}
}
