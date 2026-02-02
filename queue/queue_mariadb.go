package queue

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"utils/logging"
	"utils/queue/core"
)

type queueMariaDB struct {
	db    *sql.DB
	coreq *core.Queue
	owner string
}

type queueEntry struct {
	ID         int64
	Owner      string
	ExternalID string
	Data       sql.NullString
	CreationDT time.Time
}

func NewMariaDB(
	db *sql.DB,
	coreq *core.Queue,
	owner string,
) (
	*queueMariaDB,
	error,
) {

	if db == nil {
		return nil, errors.New("null db")
	}

	if coreq == nil {
		return nil, errors.New("null core queue")
	}

	if owner == "" {
		return nil, errors.New("empty owner")
	}

	h := &queueMariaDB{
		db:    db,
		coreq: coreq,
		owner: owner,
	}

	err := h.loadEntriesFromOwner(logging.New())
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *queueMariaDB) PushBack(
	log *logging.Logger,
	ID string,
	value string,
) error {

	l := log.New()

	err := h.insertEntry(
		l, &queueEntry{
			Owner:      h.owner,
			ExternalID: ID,
			Data: sql.NullString{
				String: value,
				Valid:  true,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error inserting entry: %w",
			err)
	}

	h.coreq.PushBack(ID, value)
	h.coreq.WakeUp()

	return nil

}

func (h *queueMariaDB) Remove(
	log *logging.Logger,
	id string,
) error {

	l := log.New()

	h.coreq.Remove(l, id)

	return h.removeEntryByExternalID(l, id)
}

func (h *queueMariaDB) Run(
	log *logging.Logger,
	f func(
		arg string,
	) bool,
) error {

	l := log.New()

	if f == nil {
		return ErrNullFunc
	}

	go h.coreq.Run(l, func(req *core.Req) bool {

		if !f(req.Value) {
			return false
		}

		// se chegou aqui, retornou true.
		// tem q remover da fila

		err := h.removeEntryByExternalID(
			l, req.ID,
		)
		if err != nil {
			l.Error("error removing entry %v: %v",
				req.ID, err)
			return false
		}

		l.Info("request %v finished successfully. removed from queue",
			req.ID)

		return true
	})

	return nil
}

// TODO finalizar a implementação
func (h *queueMariaDB) List() ([]queueEntry, error) {
	return nil, nil
}

func (h *queueMariaDB) removeEntryByExternalID(
	log *logging.Logger,
	externalID string,
) error {

	l := log.New()

	tx, err := h.db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning tx: %w", err)
	}

	cmd := `
delete from 
	queue_queue
where
	external_id = ? and
	owner = ?
`

	_, err = tx.Exec(
		cmd,
		externalID,
		h.owner,
	)
	if err != nil {
		return fmt.Errorf("error exec: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error commiting tx: %w", err)
	}

	l.Info("removed id %v", externalID)

	return nil
}

func (h *queueMariaDB) insertEntry(
	log *logging.Logger,
	entry *queueEntry,
) error {

	l := log.New()

	tx, err := h.db.Begin()
	if err != nil {
		return fmt.Errorf("error beginning tx: %w", err)
	}

	cmd := `
insert into	queue_queue(
	owner, 
	external_id, 
	data
)
values (
	?,
	?,
	?
)`

	res, err := tx.Exec(
		cmd,
		entry.Owner,
		entry.ExternalID,
		entry.Data,
	)
	if err != nil {
		return fmt.Errorf("error exec: %w", err)
	}

	l.Info("inserted %+v", entry)

	newID, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("error getting last id: %w", err)
	}

	l.Info("new id: %v", newID)

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error commiting tx: %w", err)
	}

	return nil

}

func (h *queueMariaDB) getEntriesByOwner(
	log *logging.Logger,
	owner string,
) (
	[]queueEntry,
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
	owner = ?
`

	rows, err := h.db.Query(qry, owner)
	if err != nil {
		return nil, fmt.Errorf("error on query: %w", err)
	}

	defer rows.Close()

	var entries []queueEntry

	for rows.Next() {
		var entry queueEntry
		err := rows.Scan(
			&entry.ID,
			&entry.Owner,
			&entry.ExternalID,
			&entry.Data,
			&entry.CreationDT,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning: %w", err)
		}
		entries = append(entries, entry)
	}

	l.Info("got qnty entries: %v", len(entries))

	return entries, nil

}

func (h *queueMariaDB) loadEntriesFromOwner(
	log *logging.Logger,
) error {

	l := log.New()

	entries, err := h.getEntriesByOwner(l, h.owner)
	if err != nil {
		return fmt.Errorf("error getting entries: %w", err)
	}

	for _, e := range entries {
		h.coreq.PushBack(e.ExternalID, e.Data.String)
	}
	return nil
}
