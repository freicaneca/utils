package core

import (
	"container/list"
	"fmt"
	"sync"
	"time"
	"utils/logging"
)

type Queue struct {
	queueID    string
	mutex      sync.Mutex
	qtyWorkers int

	// verification period
	periodSeconds int
	elements      *list.List
	//processFunc   func(req *Req) bool
	removedIDs map[string]struct{}
	requests   chan *Req
	wait       chan struct{}
}

// Req refers to a single element (key, value) in elements list.
type Req struct {
	ID    string
	Value string
}

func New(
	queueID string,
	qtyWorkers int,
	periodSeconds int,
	//f func(req *Req) bool,
) *Queue {

	return &Queue{
		mutex:         sync.Mutex{},
		periodSeconds: periodSeconds,
		qtyWorkers:    qtyWorkers,
		elements:      list.New(),
		//processFunc:   f,
		requests:   make(chan *Req),
		wait:       make(chan struct{}, 1),
		queueID:    queueID,
		removedIDs: make(map[string]struct{}),
	}
}

func (q *Queue) PushBack(ID string, value string) {

	//l := log.New()
	q.mutex.Lock()
	defer q.mutex.Unlock()

	rq := Req{
		Value: value,
		ID:    ID,
	}

	q.elements.PushBack(&rq)

	//q.wait <- struct{}{}

}

func (q *Queue) WakeUp() {
	q.wait <- struct{}{}
}

func (q *Queue) Run(
	log *logging.Logger,
	f func(
		req *Req,
	) bool,
) {

	l := log.New()

	for i := 0; i < q.qtyWorkers; i++ {
		go q.process(l, i, f)
	}

	for {

		q.mutex.Lock()
		e := q.elements.Front()
		q.mutex.Unlock()
		// if no element, idle wait.
		// when wait channel receives data (empty struct),
		// there will be non-nil element in Front.
		if e == nil {
			<-q.wait
			continue
		}

		// remove from queue but not from persistence.
		// it will only be removed from persistence if
		// processFunc returns true. see q.process().
		// if returns false, element will be readded to queue.
		q.removeByID(l, e.Value.(*Req).ID)

		// now sending Front() element to requests channel.
		q.requests <- e.Value.(*Req)

	}

}

func (q *Queue) process(
	log *logging.Logger,
	workerID int,
	f func(
		req *Req,
	) bool,
) {

	l := log.New()

	l.SetFrom(fmt.Sprintf("%v-worker:%v", q.queueID, workerID))

	for {

		// req contains key/value of q.elements' list
		req := <-q.requests

		//l.Info("Worker %v received request %v.", workerID, req.ID)

		// if returns true, it finished successfully.
		// if false, it did not finish ok. try again later.

		// will also check if ID is to be removed.
		// if it is, will not launch push back again.
		_, mustRemove := q.removedIDs[req.ID]

		if mustRemove {
			delete(q.removedIDs, req.ID)
		}

		if !f(req) {

			//l.Error("Request %v did not complete successfully. Will try again later.",
			//	req.ID)
			time.AfterFunc(
				time.Duration(q.periodSeconds)*time.Second,
				func() {

					if mustRemove {
						// will not push back again
						return
					}

					q.mutex.Lock()
					q.elements.PushBack(req)
					select {
					case q.wait <- struct{}{}:
					default:
					}
					q.mutex.Unlock()
				})
		}
	}

}

func (q *Queue) Remove(
	log *logging.Logger,
	ID string,
) {
	q.removedIDs[ID] = struct{}{}
	q.removeByID(log, ID)
}

func (q *Queue) removeByID(
	log *logging.Logger,
	ID string,
) {

	l := log.New()

	q.mutex.Lock()
	defer q.mutex.Unlock()

	found := false
	for e := q.elements.Front(); e != nil; e = e.Next() {
		if e.Value.(*Req).ID == ID {
			q.elements.Remove(e)
			found = true
			//l.Info("removed element with ID %q", ID)
			break
		}
	}

	if !found {
		l.Warn("Did not find element with ID %q", ID)
	}
}
