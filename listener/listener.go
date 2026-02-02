package listener

import (
	"time"
	"utils/logging"
	"utils/messenger"
	"utils/workerpool"
)

type BrokerListener struct {
	topic        string
	inbox        string
	numOfWorkers uint
	worker       workerpool.WorkerFunc
	messenger    messenger.Client
}

func NewBrokerListener(
	messenger messenger.Client,
	topic string,
	inbox string,
	numOfWorkers uint,
	worker workerpool.WorkerFunc,
) *BrokerListener {
	return &BrokerListener{
		topic:        topic,
		inbox:        inbox,
		messenger:    messenger,
		numOfWorkers: numOfWorkers,
		worker:       worker,
	}
}

func (bl *BrokerListener) Listen(log *logging.Logger) {
	l := log.New()
	var reader messenger.Reader
	var err error

	for {
		reader, err = bl.messenger.NewReader(
			bl.topic,
			bl.inbox,
			messenger.SharedInbox,
			true,
		)
		if err != nil {
			l.Error("Error setting up messenger: %v", err)
			l.Error("Will try again soon.")
			time.Sleep(5 * time.Second)
			continue
		}

		break
	}

	wp := workerpool.New(bl.numOfWorkers, bl.worker)
	wp.Start(l)

	defer reader.Close(l)

	l.Info("Now listening to topic %q.", bl.topic)

	for {
		// Get the mensagem
		brokerMsg, err := reader.Peek(2 * time.Minute)
		if err != nil {
			_, isTimeout := err.(*messenger.TimeoutError)
			if !isTimeout {
				l.Debug("Error getting message from router: %v", err)
				brokerMsg.GiveBack()
			}

			continue
		}
		wp.Feed(brokerMsg)
	}
}
