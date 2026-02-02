package messenger

import (
	"time"
	"utils/logging"
)

type Message interface {
	WriteToModel(
		model any,
	) error

	Payload() []byte

	Received()

	GiveBack()
}

type Reader interface {
	Get(
		model any,
		timeout time.Duration,
	) error

	Peek(
		timeout time.Duration,
	) (
		Message,
		error,
	)

	Close(
		log *logging.Logger,
	)
}

type Client interface {
	NewReader(
		from string,
		inboxName string,
		inboxType InboxType,
		ignorePreviousMessages bool,
	) (
		Reader,
		error,
	)

	Send(
		to string,
		model any,
	) error

	Close()
}
