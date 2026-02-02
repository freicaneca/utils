package queue

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"
	"utils/logging"
	"utils/queue/core"
	"utils/utils/testutils"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

func TestQueueMariaDB(t *testing.T) {

	l := logging.New()

	db, err := sql.Open(
		"mysql",
		"root:baba@tcp(localhost:3306)/test?parseTime=true",
		//"filipe:#UmniDb123@tcp(localhost:3306)/test?parseTime=true",
	)
	if err != nil {
		t.Fatalf("error opening db: %v", err)
	}

	coreq := core.New("test", 1, 1)

	owner := "duduq"

	h, err := NewMariaDB(db, coreq, owner)
	testutils.AssertBool(t, err == nil, true)

	t.Run("load entries by owner at initialization", func(t *testing.T) {

		// first, initialize with no entries

		_, err := NewMariaDB(db, coreq, "lalalae")
		testutils.AssertBool(t, err == nil, true)

		// now, let's insert some data and load a new queue
		// for a new owner

		qty := 3
		for i := range qty {
			value := fmt.Sprintf("test %d", i)
			uu := uuid.NewString()
			h.insertEntry(
				l, &queueEntry{
					Owner:      owner,
					ExternalID: uu,
					Data: sql.NullString{
						String: value,
						Valid:  true,
					},
				},
			)

			defer h.removeEntryByExternalID(l, uu)
		}

		_, err = NewMariaDB(db, coreq, owner)
		testutils.AssertBool(t, err == nil, true)

		// entries, err := newH.List()
		// if err != nil {
		// 	t.Fatalf("error selecting entries from %s: %v", owner, err)
		// }

		// testutils.AssertInt(t, len(entries), qty)

	})

	t.Run("tempo de vida", func(t *testing.T) {

		type data struct {
			CreationDT time.Time `json:"creation_date_time"`
		}

		maxAge := 3 * time.Second

		err := h.Run(
			l, func(arg string) bool {

				d := data{}

				err := json.Unmarshal([]byte(arg), &d)
				if err != nil {
					t.Fatalf("error umarshalling: %v", err)
				}

				if time.Since(d.CreationDT) > maxAge {
					l.Info("max age reached! removing...")
					return true
				}

				l.Info("max age not reached. continuing...")
				return false

			},
		)
		if err != nil {
			t.Fatalf("error running: %v", err)
		}

		d := data{
			CreationDT: time.Now(),
		}

		dB, _ := json.Marshal(d)

		id := uuid.NewString()

		err = h.PushBack(l, id, string(dB))
		if err != nil {
			t.Fatalf("error pushing back: %v", err)
		}

		time.Sleep(5 * time.Second)

		// depois de 5s, esperamos q o item tenha sido removido.
		// vamos dar um get e esperamos obter erro de not found

		_, err = getEntryByExternalID(
			l, db, id,
		)
		testutils.AssertBool(t, err != nil, true)
		testutils.AssertBool(
			t,
			errors.Is(err, sql.ErrNoRows),
			true,
		)

	})

	t.Run("get entries by owner", func(t *testing.T) {

		qty := 3
		for i := range qty {
			value := fmt.Sprintf("test %d", i)
			uu := uuid.NewString()
			h.insertEntry(
				l, &queueEntry{
					Owner:      owner,
					ExternalID: uu,
					Data: sql.NullString{
						String: value,
						Valid:  true,
					},
				},
			)

			defer h.removeEntryByExternalID(l, uu)
		}

		entries, err := h.getEntriesByOwner(l, owner)
		if err != nil {
			t.Fatalf("error selecting entries from %s: %v", owner, err)
		}

		testutils.AssertInt(t, len(entries), qty)

		for _, e := range entries {
			testutils.AssertBool(t, e.Owner == owner, true)
		}

	})
}

func getEntryByExternalID(
	log *logging.Logger,
	db *sql.DB,
	externalID string,
) (
	*queueEntry,
	error,
) {

	l := log.New()

	qry := `
select 
	id, 
	owner,
	external_id,
	data,
	creation_date_time
from 
	queue_queue
where
	external_id = ?
`

	row := db.QueryRow(
		qry, externalID,
	)

	out := queueEntry{}

	err := row.Scan(
		&out.ID,
		&out.Owner,
		&out.ExternalID,
		&out.Data,
		&out.CreationDT,
	)
	if err != nil {
		return nil, fmt.Errorf("error scanning: %w", err)
	}

	l.Info("got entry: %+v", out)

	return &out, nil

}
