package broadcaster

import (
	"time"

	"github.com/pkg/errors"
	tele "gopkg.in/telebot.v3"
)

type Broadcaster struct {
	bot *tele.Bot

	BaseWait     time.Duration
	AddFloodWait time.Duration // adds to RetryAfter field

	errHandler ErrorHandler
}

func NewBroadcaster(
	bot *tele.Bot,
	baseWait time.Duration,
	floodWait time.Duration,
	errHandler ErrorHandler,
) *Broadcaster {
	if baseWait == 0 {
		baseWait = DefaultWait
	}

	if errHandler == nil {
		errHandler = func(_ *Error) (next bool) {
			return false // returning last error and stop broadcast
		}
	}

	return &Broadcaster{
		bot:          bot,
		BaseWait:     baseWait,
		AddFloodWait: floodWait,
		errHandler:   errHandler,
	}
}

type ErrorHandler func(err *Error) (next bool)

var DefaultWait = time.Millisecond * 60 // ~20 messages per sec. Limit 30 per sec.

// Send executes broadcast. It's thread-blocked function.
// Recipients will be taken from the channel. This solution allows you to:
//		1. Do not get all the chats from the database, but send them one by one or in batches
//		2. Convenient filtering in "real" time
// Sending object (what) is tele.Sendable for exclude errors about the wrong type of object being sent.
// If you need send text use TextSender.
func (b *Broadcaster) Send(what tele.Sendable, chats chan tele.Recipient, opts *tele.SendOptions) (successfully int, err error) {
	if opts == nil {
		opts = new(tele.SendOptions)
	}

	for chat := range chats {
		time.Sleep(b.BaseWait)

		_, err = b.bot.Send(chat, what, opts)
		if err == nil {
			successfully++
			continue
		}

		var floodError *tele.FloodError
		if errors.As(err, &floodError) {
			retryAfter := time.Second * time.Duration(floodError.RetryAfter)
			time.Sleep(b.AddFloodWait + retryAfter)
			successfully++
			continue
		}

		if !b.errHandler(&Error{In: chat, Err: err}) {
			return
		}
	}

	return
}

func (b *Broadcaster) SetErrHandler(errHandler ErrorHandler) {
	if errHandler == nil {
		return
	}

	b.errHandler = errHandler
}

func (b *Broadcaster) Copy() *Broadcaster {
	bc := new(Broadcaster)
	*bc = *b
	return bc
}
