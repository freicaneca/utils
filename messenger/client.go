package messenger

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"utils/logging"

	"github.com/apache/pulsar-client-go/pulsar"
)

// TimeoutError the requested operation has reached a timeout
type TimeoutError struct {
	Err error
}

func (te *TimeoutError) Error() string {
	return fmt.Sprintf("timeout error: %#v", te.Err)
}

// Message represents a message
type pulsarMessage struct {
	consumer pulsar.Consumer
	message  pulsar.Message
}

func newMessage(
	message pulsar.Message,
	consumer pulsar.Consumer,
) *pulsarMessage {
	return &pulsarMessage{
		consumer: consumer,
		message:  message,
	}
}

// WriteToModel writes message's content to a model
func (m *pulsarMessage) WriteToModel(
	model any,
) error {
	err := json.Unmarshal(m.message.Payload(), model)
	if err != nil {
		return fmt.Errorf("failed to decode message: %#v", err)
	}
	return nil
}

// Payload returns message's payload
func (m *pulsarMessage) Payload() []byte {
	return m.message.Payload()
}

// Received marks a message as received
func (m *pulsarMessage) Received() {
	m.consumer.Ack(m.message)
}

// GiveBack gives the message back, so it can be readed again
func (m *pulsarMessage) GiveBack() {
	m.consumer.Nack(m.message)
}

// InboxType determines the message reception behavior (exclusive, shared and so on)
type InboxType int

const (
	// ExclusiveInbox the reader will receive its exlusive copy of the message in his inbox,
	// and this inbox can't be shared with other readers
	ExclusiveInbox InboxType = iota

	// SharedInbox the reader will receive a message in his inbox, and the messages will be
	// dispatched according to a round-robin rotation between readers sharing the same inbox
	SharedInbox

	// FailoverInbox only one reader (the first to arrive) can read messages from this inbox,
	// and the others have to wait their turn in line until the first abandon it
	FailoverInbox

	// KeyShared subscription mode, multiple consumer will be able to use the same
	// subscription and all messages with the same key will be dispatched to only one consumer
	//KeyShared
)

// Reader represents a message reader
type pulsarReader struct {
	consumer pulsar.Consumer
}

// Get receives a message, writes its content to a model and mark it as received
func (r *pulsarReader) Get(model interface{}, timeout time.Duration) error {
	msg, err := r.Peek(timeout)
	if err != nil {
		return fmt.Errorf("failed to get message: %v", err)
	}
	if err := msg.WriteToModel(model); err != nil {
		return fmt.Errorf("failed to write message to model: %#v", err)
	}
	msg.Received()
	return nil
}

// Peek receives a message, but doesn't mark it as received
func (r *pulsarReader) Peek(
	timeout time.Duration,
) (
	Message,
	error,
) {
	ctx, canc := context.WithTimeout(context.Background(), timeout)
	defer canc()
	msg, err := r.consumer.Receive(ctx)
	if err != nil {
		if ctx.Err() != nil {
			return nil, &TimeoutError{Err: err}
		}
		return nil, fmt.Errorf("failed to consume message: %#v", err)
	}
	return newMessage(msg, r.consumer), nil
}

// Close finishes the reader
func (r *pulsarReader) Close(log *logging.Logger) {
	l := log.New()
	if err := r.consumer.Unsubscribe(); err != nil {
		l.Error("failed to unsubscribe: %#v", err)
	}
	r.consumer.Close()
}

// Client represents a message broker client
type pulsarClient struct {
	client pulsar.Client
}

// NewClient creates an instance of a messenger client
func NewClient(
	messengerURL string,
	connectionTimeout time.Duration,
	operationTimeout time.Duration,
) (*pulsarClient, error) {
	client, err := pulsar.NewClient(pulsar.ClientOptions{
		URL:               messengerURL,
		ConnectionTimeout: connectionTimeout,
		OperationTimeout:  operationTimeout,
	})
	if err != nil {
		return nil, fmt.Errorf("could not instantiate Pulsar client: %#v", err)
	}
	return &pulsarClient{
		client: client,
	}, nil
}

// Send posts a model based message to message broker
func (m *pulsarClient) Send(to string, model interface{}) error {
	producer, err := m.client.CreateProducer(pulsar.ProducerOptions{
		Topic: to,
	})
	if err != nil {
		return fmt.Errorf("failed to create producer: %#v", err)
	}
	defer producer.Close()
	payload, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("failed to encode message: %#v", err)
	}
	_, err = producer.Send(context.Background(), &pulsar.ProducerMessage{
		Payload: payload,
	})
	if err != nil {
		return fmt.Errorf("failed to publish message: %#v", err)
	}
	return nil
}

// NewReader creates a message reader based on the messenger client
func (m *pulsarClient) NewReader(
	from string,
	inboxName string,
	inboxType InboxType,
	ignorePreviousMessages bool,
) (
	Reader,
	error,
) {
	var stype pulsar.SubscriptionType
	switch inboxType {
	case ExclusiveInbox:
		stype = pulsar.Exclusive
	case SharedInbox:
		stype = pulsar.Shared
	case FailoverInbox:
		stype = pulsar.Failover
	}
	initPosition := pulsar.SubscriptionPositionEarliest
	if ignorePreviousMessages {
		initPosition = pulsar.SubscriptionPositionLatest
	}
	consumer, err := m.client.Subscribe(pulsar.ConsumerOptions{
		Topic:                       from,
		SubscriptionName:            inboxName,
		Type:                        stype,
		SubscriptionInitialPosition: initPosition,
		NackRedeliveryDelay:         1 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &pulsarReader{consumer: consumer}, nil
}

// Close closes client connection
func (m *pulsarClient) Close() {
	m.client.Close()
}
