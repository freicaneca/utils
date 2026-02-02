package svcregistry

import (
	"context"
	"errors"
	"utils/logging"
)

// getter de registro de serviço bem simples.
type simpleGetter struct {
	// key: nome do serviço
	// value: endereço do serviço
	addresses map[string]string
}

func NewSimpleGetter(
	addresses map[string]string,
) (
	*simpleGetter,
	error,
) {

	if addresses == nil {
		return nil, errors.New("null addresses")
	}

	return &simpleGetter{
		addresses: addresses,
	}, nil

}

func (h *simpleGetter) GetServiceAddress(
	log *logging.Logger,
	ctx context.Context,
	service string,
) (
	string,
	error,
) {

	return h.addresses[service], nil

}

func (h *simpleGetter) Reset(
	log *logging.Logger,
	service string,
) {
}
