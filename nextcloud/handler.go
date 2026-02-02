package nextcloud

import (
	"context"
	"errors"
	"utils/logging"
)

type Handler interface {
	RegisterUser(
		log *logging.Logger,
		ctx context.Context,
		userID string,
		password string,
	) error

	RemoveUser(
		log *logging.Logger,
		ctx context.Context,
		userID string,
	) error
}

var (
	ErrBadRequest        = errors.New("bad request")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrBadPassword       = errors.New("bad password")
	ErrInternal          = errors.New("internal err")
)
