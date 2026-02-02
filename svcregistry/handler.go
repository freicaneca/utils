package svcregistry

import (
	"context"
	"time"
	"utils/logging"
)

type Handler interface {
	RegisterInstance(
		log *logging.Logger,
		ctx context.Context,
		instanceID string,
		serviceName string,
		address string,
		port int,
		tags []string,
		checkInterval time.Duration,
	) error

	DeregisterInstance(
		log *logging.Logger,
		ctx context.Context,
		instanceID string,
	) error

	SendHeartbeat(
		log *logging.Logger,
		ctx context.Context,
		instanceID string,
		serviceName string,
	) error
}
