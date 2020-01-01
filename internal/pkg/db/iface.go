package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
)

//ErrBadToken indicates that an invalid pagination token has been provided
var ErrBadToken = errors.New("bad pagination token")

//DB is the minimal database interface to back the app
//go:generate counterfeiter . DB
type DB interface {
	EnsureDBConfig(ctx context.Context) error
	EnsureDevicesTable(ctx context.Context) error
	EnsureReadingsTable(ctx context.Context) error

	RegisterDevice(ctx context.Context, name string, config string) (string, error)
	ConfigureDevice(ctx context.Context, uuid string, config string) error
	GetDeviceByUUID(ctx context.Context, uuid string) (Device, error)

	InsertReading(ctx context.Context, deviceUUID string, tCelcius json.Number, rh json.Number) (string, error)

	//TODO QueryReadings w/ pagination
}

//DBAgent implements DB using a *sql.DB
type DBAgent struct {
	db     *sql.DB
	uuider UUIDer
}

//NewDBAgent create a new DBAgent
func NewDBAgent(db *sql.DB) *DBAgent {
	return &DBAgent{
		db:     db,
		uuider: UUIDGenerator{},
	}
}

//EnsureDBConfig EnsureDBConfig
func (a *DBAgent) EnsureDBConfig(ctx context.Context) error {
	_, err := a.db.ExecContext(ctx, `CREATE EXTENSION IF NOT EXISTS "pgcrypto"`)
	if err != nil {
		return fmt.Errorf("failed to ensure pgcrypto extension: %w", err)
	}
	return nil
}

//EnsureDatabase builds out all the tables in order
func EnsureDatabase(ctx context.Context, db DB) error {
	err := db.EnsureDBConfig(ctx)
	if err != nil {
		return err
	}
	err = db.EnsureDevicesTable(ctx)
	if err != nil {
		return err
	}
	err = db.EnsureReadingsTable(ctx)
	if err != nil {
		return err
	}
	return nil
}
