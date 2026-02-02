package svcregistry

import (
	"context"
	"errors"
	"fmt"
	"time"
	"utils/logging"

	"github.com/hashicorp/consul/api"
)

type consulHandler struct {
	consulAddress string
	consulPort    uint
	client        *api.Client
}

func NewConsulHandler(
	consulAddress string,
	consulPort uint,
) (
	*consulHandler,
	error,
) {

	if len(consulAddress) == 0 {
		return nil, errors.New("empty server address")
	}

	if consulPort == 0 {
		return nil, errors.New("invalid server port")
	}

	client, err := api.NewClient(
		&api.Config{
			Address: consulAddress + ":" +
				fmt.Sprintf("%v", consulPort),
		},
	)
	if err != nil {
		return nil, err
	}

	return &consulHandler{
		consulAddress: consulAddress,
		consulPort:    consulPort,
		client:        client,
	}, nil

}

func (h *consulHandler) RegisterInstance(
	log *logging.Logger,
	ctx context.Context,
	serviceName string,
	instanceID string,
	address string,
	port int,
	tags []string,
	checkInterval time.Duration,
) error {

	l := log.New()

	registration := api.AgentServiceRegistration{
		ID:      instanceID,
		Name:    serviceName,
		Address: address,
		Port:    port,
		Tags:    tags,
	}

	if checkInterval > 0 {
		registration.Check = &api.AgentServiceCheck{
			CheckID:                        fmt.Sprintf("%s_%s", serviceName, instanceID),
			TTL:                            checkInterval.String(),
			DeregisterCriticalServiceAfter: (100 * checkInterval).String(),
		}
	}
	err := h.client.Agent().ServiceRegister(&registration)
	if err != nil {
		l.Error("error registering service: %+v",
			err)
		return fmt.Errorf("error registering service: %v",
			err)
	}

	return nil
}

func (h *consulHandler) DeregisterInstance(
	log *logging.Logger,
	ctx context.Context,
	instanceID string,
) error {

	l := log.New()

	err := h.client.Agent().ServiceDeregister(instanceID)
	if err != nil {
		l.Error("failed to deregister instance: %v", err)
		return fmt.Errorf("failed to deregister service: %w", err)
	}

	l.Debug("instance deregistered successfully: %s", instanceID)
	return nil
}

func (h *consulHandler) SendHeartbeat(
	log *logging.Logger,
	ctx context.Context,
	serviceName string,
	instanceID string,
) error {

	l := log.New()

	checkID := fmt.Sprintf(
		"%s_%s", serviceName, instanceID)
	err := h.client.Agent().PassTTL(
		checkID, "last heartbeat at: "+time.Now().Format(time.RFC3339))
	if err != nil {
		l.Error("error sending heartbeat: %+v", err)
		return fmt.Errorf(
			"failed to send heartbeat for service %s: %v",
			serviceName, err)
	}

	return nil
}
