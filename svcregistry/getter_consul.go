package svcregistry

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
	"utils/logging"

	"github.com/hashicorp/consul/api"
)

type consulGetter struct {
	serverAddress string
	serverPort    uint
	client        *api.Client
	cacheDuration time.Duration
	cacheLock     *sync.Mutex
	// key: service name
	cachedAddresses map[string]cacheControl
}

func NewConsulGetter(
	serverAddress string,
	serverPort uint,
	cacheDuration time.Duration,
) (
	*consulGetter,
	error,
) {

	if len(serverAddress) == 0 {
		return nil, errors.New("empty server address")
	}

	if serverPort == 0 {
		return nil, errors.New("invalid server port")
	}

	client, err := api.NewClient(&api.Config{
		Address: serverAddress + ":" +
			fmt.Sprintf("%v", serverPort),
	})
	if err != nil {
		return nil, err
	}

	return &consulGetter{
		serverAddress:   serverAddress,
		serverPort:      serverPort,
		client:          client,
		cacheDuration:   cacheDuration,
		cachedAddresses: map[string]cacheControl{},
		cacheLock:       &sync.Mutex{},
	}, nil

}

func (h *consulGetter) GetServiceAddress(
	log *logging.Logger,
	ctx context.Context,
	service string,
) (
	string,
	error,
) {

	l := log.New()

	nowDT := time.Now()

	cacheCtrl, ok := h.cachedAddresses[service]
	if (!ok || len(cacheCtrl.addresses) == 0) ||
		(len(cacheCtrl.addresses) > 0 &&
			nowDT.After(cacheCtrl.creationDT.Add(h.cacheDuration))) {

		l.Debug("service %q not found within cache",
			service)
		l.Debug("or cache has expired")
		l.Debug("will check with consul and update cache")

		insts, err := h.getServiceInstancesFromConsul(
			l, ctx, service,
		)
		if err != nil {
			l.Error("error getting svcs from consul: %+v",
				err)
			return "", err
		}

		h.cacheLock.Lock()
		addrs := make([]string, 0, len(insts))

		nowDT := time.Now()
		addrs = append(addrs, insts...)
		// for _, inst := range insts {
		// 	addrs = append(addrs, inst)
		// }

		h.cachedAddresses[service] = cacheControl{
			creationDT: nowDT,
			cursor:     0,
			addresses:  addrs,
		}
		h.cacheLock.Unlock()

		cacheCtrl = h.cachedAddresses[service]

	} else {
		l.Debug("found service %v in cache and it's valid",
			service)
	}

	// pegando uma instancia aleatoria
	if len(cacheCtrl.addresses) == 0 {
		l.Error("service not found: %v", service)
		return "", ErrServiceNotFound
	}

	inst := h.nextAddress(service)
	//h.cachedAddresses[service] = cacheCtrl

	l.Debug("got instance addr: %v",
		inst)

	return inst, nil
}

// busca o próximo endereço para o dado serviço
// e atualiza o cursor:
// política round robin.
func (h *consulGetter) nextAddress(
	service string,
) string {

	cacheCtrl, ok := h.cachedAddresses[service]
	if !ok {
		return ""
	}

	lenAddrs := len(cacheCtrl.addresses)

	if lenAddrs == 0 {
		cacheCtrl.cursor = 0
		return ""
	}

	out := cacheCtrl.addresses[cacheCtrl.cursor]

	cacheCtrl.cursor++

	if cacheCtrl.cursor >= lenAddrs {
		cacheCtrl.cursor = 0
	}

	h.cacheLock.Lock()
	h.cachedAddresses[service] = cacheCtrl
	h.cacheLock.Unlock()

	return out
}

func (h *consulGetter) getServiceInstancesFromConsul(
	log *logging.Logger,
	ctx context.Context,
	service string,
) (
	[]string,
	error,
) {

	l := log.New()

	opts := &api.QueryOptions{}
	opts = opts.WithContext(ctx)

	cSvcs, _, err := h.client.Health().Service(
		service, "", true, opts)
	if err != nil {
		l.Error("failed to get service %q: %v", service, err)
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	if len(service) == 0 {
		l.Debug("service %q not found", service)
		return []string{}, nil
	}

	out := make(
		[]string, 0, len(cSvcs))
	for _, s := range cSvcs {
		out = append(out, s.Service.Address+":"+
			strconv.Itoa(s.Service.Port))
		// out = append(out, &ServiceAddress{
		// 	Address: s.Service.Address,
		// 	Port:    s.Service.Port,
		// })
	}
	return out, nil

}

func (h *consulGetter) Reset(
	log *logging.Logger,
	service string,
) {

	l := log.New()

	h.cacheLock.Lock()
	defer h.cacheLock.Unlock()

	// removendo o endereço usado na última tentativa
	// e resetando o cursor.

	cCtrl, ok := h.cachedAddresses[service]
	if !ok {
		return
	}

	lenAddrs := len(cCtrl.addresses)

	// se tiver 0 ou 1 endereço, remove do cache.
	if lenAddrs <= 1 {
		delete(h.cachedAddresses, service)
		return
	}

	// se tá aqui, é pq o cache tem > 1 endereço.
	// vamos remover o último utilizado, q pode ser
	// encontrado a partir do cursor:
	// tem q subtrair 1 do cursor (com cuidado pra nao dar
	// negativo.)
	newAddrs := make([]string, 0, lenAddrs-1)

	lastPos := cCtrl.cursor - 1
	if lastPos < 0 {
		lastPos = lenAddrs - 1
	}

	for i, addr := range cCtrl.addresses {
		if i == lastPos {
			continue
		}

		newAddrs = append(newAddrs, addr)
	}

	h.cachedAddresses[service] = cacheControl{
		creationDT: cCtrl.creationDT,
		cursor:     0,
		addresses:  newAddrs,
	}

	l.Info("cache has been reset for %q", service)

}

type cacheControl struct {
	// data de criacao do cache
	creationDT time.Time

	// cursor atual da lista dos endereços.
	// a cada chamada de GetServiceAddress, o cursor
	// é incrementado (round robin).
	cursor int

	// lista de endereços de um dado serviço
	addresses []string
}
