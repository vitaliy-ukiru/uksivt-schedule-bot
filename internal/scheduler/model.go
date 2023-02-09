package scheduler

import (
	"time"
)

type FlagSet uint16

const (
	NextDay FlagSet = 1 << iota
	Full
	OnlyIfHaveReplaces
	ReplacesAlways
	FullOnlyIfReplaces
)

type CronJob struct {
	ID     int64     `json:"id" mapstructure:"-"`
	ChatID int64     `json:"chat_id" mapstructure:"-"`
	Title  string    `json:"title" mapstructure:"title,omitempty"`
	At     time.Time `json:"at" mapstructure:"send_at,omitempty"`
	Flags  FlagSet   `json:"flags,omitempty" mapstructure:"flags,omitempty"`
}

func (f FlagSet) Has(b FlagSet) bool {
	return (f & b) != 0
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

func (f FlagSet) With(b FlagSet) FlagSet {
	return f | b
}

func (f FlagSet) Without(b FlagSet) FlagSet {
	return f & (^b)
}
