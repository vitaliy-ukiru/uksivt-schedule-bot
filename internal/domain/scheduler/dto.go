package scheduler

import "time"

type CreateJobDTO struct {
	Title  string    `json:"title"`
	At     time.Time `json:"at"`
	Flags  uint16    `json:"flags,omitempty"`
	ChatID int64     `json:"chat_id,omitempty"`
}

type UpdateDTO struct {
	NewTime  *time.Time `json:"new_time,omitempty"`
	NewFlags *uint16    `json:"new_flags,omitempty"`
}
