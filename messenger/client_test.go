package messenger

import (
	"testing"
	"time"
	"utils/utils/testutils"
)

func TestClient(t *testing.T) {

	pulsar, err := NewClient(
		"pulsar://relaunch.umnicorn.com:6650",
		10*time.Second,
		5*time.Second,
	)
	testutils.AssertError(t, err, nil)

	err = pulsar.Send("baba", struct {
		Baba string
	}{
		Baba: "bobo",
	})
	testutils.AssertError(t, err, nil)
}
