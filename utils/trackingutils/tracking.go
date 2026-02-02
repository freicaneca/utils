package trackingutils

import (
	"fmt"
	"sync/atomic"
	"time"
)

var GlobalTrackingNumber TrackingNumber

type TrackingNumber struct {
	number uint64
}

func (tn *TrackingNumber) Next() string {
	return fmt.Sprintf("%v.%v", time.Now().UnixNano(), atomic.AddUint64(&tn.number, 1))
}
