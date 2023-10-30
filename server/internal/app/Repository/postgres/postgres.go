package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

const TIMESTAMP_FORMAT = time.RFC3339

var (
	ErrWrongScannedValue    = errors.New("scanned value isn't valid")
	ErrPragmaIntegrityNotOk = errors.New("pragma integrity_check is not ok")
)

type InstanceInterface interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Instance struct {
	connection *sql.DB
	// Announce   *AnnounceModel
	// Media      *MediaModel
	// Phone      *PhoneModel
}

func MustOpen(connStr string) *Instance {
	inst, err := Open(connStr)
	if err != nil {
		panic(err)
	}

	return inst
}

func Open(connStr string) (*Instance, error) {
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	inst := &Instance{
		connection: db,
		// Announce:   &AnnounceModel{db: db},
		// Media:      &MediaModel{db: db},
		// Phone:      &PhoneModel{db: db},
	}

	return inst, err
}

func (db *Instance) Connection() *sql.DB {
	return db.connection
}
