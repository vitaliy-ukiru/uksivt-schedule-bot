package chat

import "errors"

var (
	ErrChatDeleted  = errors.New("chat was deleted")
	ErrChatNotFound = errors.New("chat not found")
	ErrNotModified  = errors.New("not modified")
)
