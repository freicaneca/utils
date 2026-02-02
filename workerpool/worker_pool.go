package workerpool

import (
	"fmt"
	"utils/logging"
)

type WorkerFunc func(
	log *logging.Logger, id string, in chan any)

type WorkersPool struct {
	numOfWorkers uint
	in           chan any
	worker       WorkerFunc
}

func New(
	numOfWorkers uint,
	worker WorkerFunc,
) *WorkersPool {
	return &WorkersPool{
		in:           make(chan any),
		numOfWorkers: numOfWorkers,
		worker:       worker,
	}
}

func (wp *WorkersPool) Start(log *logging.Logger) {
	for i := 0; i < int(wp.numOfWorkers); i++ {
		id := fmt.Sprintf("worker-%v", i)
		go wp.worker(log, id, wp.in)
	}
}

func (wp *WorkersPool) Feed(payload any) {
	wp.in <- payload
}

func (wp *WorkersPool) AsyncFeed(payload any) bool {
	select {
	case wp.in <- payload:
		return true
	default:
		return false
	}
}
