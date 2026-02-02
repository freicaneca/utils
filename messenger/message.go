package messenger

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"utils/logging"
	"utils/utils/contextutils"
)

type CommonMessage struct {
	TrackingID       string `json:"tracking_id"`
	Origin           string `json:"origin"`
	CreationDateTime string `json:"creation_date_time"`
	Data             string `json:"data"`
}

// Publishes a CommonMessage whose "data" field is
// "data" serialized.
func PublishMessage(
	log *logging.Logger,
	ctx context.Context,
	client Client,
	origin string,
	topic string,
	data any,
) error {

	l := log.New()

	if data == nil {
		return fmt.Errorf("null data")
	}

	if topic == "" {
		return fmt.Errorf("empty topic")
	}

	dataB, _ := json.Marshal(data)

	nowDT := time.Now()

	msg := CommonMessage{
		TrackingID: contextutils.GetContextValue(
			ctx, contextutils.ContextKeyReqTracking).(string),
		Origin:           origin,
		CreationDateTime: nowDT.Format(time.RFC3339Nano),
		Data:             string(dataB),
	}

	err := client.Send(topic, msg)
	if err != nil {
		return fmt.Errorf("%w: err publishing msg",
			err)
	}

	l.Info("published msg to topic %q", topic)

	return nil
}
