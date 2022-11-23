package scheduler

import (
	"time"

	"github.com/vitaliy-ukiru/uksivt-schedule-bot/internal/domain/chat"
)

type FlagSet uint16

const (
	_ FlagSet = iota << 1
	Full
	OnlyIfHaveReplaces
	ReplacesAlways
	FullOnlyIfReplaces
	NextDay
)

type CronJob struct {
	ID    int64
	At    time.Time
	Flags FlagSet
	Chat  *chat.Chat
}

func (f FlagSet) Has(b FlagSet) bool {
	return (f & b) > 0
}

func (f FlagSet) HasAny(other ...FlagSet) bool {
	for _, b := range other {
		if f.Has(b) {
			return true
		}
	}
	return false
}

func (f *FlagSet) Add(b FlagSet) {
	*f |= b
}

func (f *FlagSet) Unset(b FlagSet) {
	*f &= ^b
}

type CronJobBase struct {
	ID     int64
	At     time.Time
	Flags  FlagSet
	ChatID int64
}
