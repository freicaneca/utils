package queue

import (
	"errors"
	"utils/logging"
)

type Handler interface {
	PushBack(
		log *logging.Logger,
		ID string,
		value string,
	) error

	Run(
		log *logging.Logger,
		f func(
			arg string,
		) bool,
	) error

	Remove(
		log *logging.Logger,
		ID string,
	) error
}

var ErrNullFunc error = errors.New("null function")
