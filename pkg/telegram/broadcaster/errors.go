package broadcaster

import (
	"fmt"

	tele "gopkg.in/telebot.v3"
)

type Error struct {
	In  tele.Recipient
	Err error
}

func (e Error) Error() string {
	return fmt.Sprintf("brodcaster: chat(%s), error: %s", e.In.Recipient(), e.Err)
}

func (e Error) Cause() error  { return e.Err }
func (e Error) Unwrap() error { return e.Err }

// ErrorProxying proxying errors to channel ch.
func ErrorProxying(ch chan *Error) ErrorHandler {
	return func(err *Error) bool {
		ch <- err
		return true
	}
}
