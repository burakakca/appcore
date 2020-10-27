package models

import (
	"database/sql"
	"net/url"

	"github.com/go-fed/activity/streams/vocab"
	"github.com/go-fed/apcore/util"
)

var _ Model = &LocalData{}

// LocalData is a Model that provides additional database methods for
// ActivityStreams data generated by this instance.
type LocalData struct {
	localCreate *sql.Stmt
	localUpdate *sql.Stmt
	localDelete *sql.Stmt
}

func (f *LocalData) Prepare(db *sql.DB, s SqlDialect) error {
	var err error
	f.localCreate, err = db.Prepare(s.LocalCreate())
	if err != nil {
		return err
	}
	f.localUpdate, err = db.Prepare(s.LocalUpdate())
	if err != nil {
		return err
	}
	f.localDelete, err = db.Prepare(s.LocalDelete())
	if err != nil {
		return err
	}
	return nil
}

func (f *LocalData) CreateTable(t *sql.Tx, s SqlDialect) error {
	_, err := t.Exec(s.CreateLocalDataTable())
	return err
}

func (f *LocalData) Close() {
	f.localCreate.Close()
	f.localUpdate.Close()
	f.localDelete.Close()
}

// Create inserts the local data into the table.
func (f *LocalData) Create(c util.Context, tx *sql.Tx, v vocab.Type) error {
	_, err := tx.Stmt(f.localCreate).ExecContext(c, ActivityStreams{v})
	return err
}

// Update replaces the local data for the specified IRI.
func (f *LocalData) Update(c util.Context, tx *sql.Tx, localIDIRI *url.URL, v vocab.Type) error {
	_, err := tx.Stmt(f.localUpdate).ExecContext(c, localIDIRI.String(), ActivityStreams{v})
	return err
}

// Delete removes the local data with the specified IRI.
func (f *LocalData) Delete(c util.Context, tx *sql.Tx, localIDIRI *url.URL) error {
	_, err := tx.Stmt(f.localDelete).ExecContext(c, localIDIRI.String())
	return err
}
