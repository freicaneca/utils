package svcregistry

import (
	"context"
	"utils/logging"
)

type Getter interface {
	// retorna o endereço de uma instância do serviço "service".
	// pode usar cache interno (se houver).
	GetServiceAddress(
		log *logging.Logger,
		ctx context.Context,
		service string,
	) (
		string,
		error,
	)

	// remove o último endereço utilizado
	// de service.
	// isso vai forçar que um outro endereço seja utilizado
	// na próxima chamada de GetServiceAddress
	Reset(
		log *logging.Logger,
		service string,
	)
}

// type ServiceAddress struct {
// 	Address string
// 	Port    int
// }
