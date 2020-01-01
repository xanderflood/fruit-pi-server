package db

import (
	"context"
	"encoding/json"
	"fmt"
)

//EnsureReadingsTable EnsureReadingsTable
func (a *DBAgent) EnsureReadingsTable(ctx context.Context) error {
	_, err := a.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS "readings"
(	"uuid" UUID DEFAULT gen_random_uuid(),
	"created_at" timestamp NOT NULL,
	"modified_at" timestamp NOT NULL,
	"deleted_at" timestamp,

	"device_uuid" UUID REFERENCES devices(uuid) NOT NULL,
	"temperature_celcius" varchar NOT NULL,
    "relative_humidity" varchar NOT NULL,
	PRIMARY KEY ("uuid")
)`)
	if err != nil {
		return fmt.Errorf("failed to ensure accounts table: %w", err)
	}

	_, err = a.db.ExecContext(ctx, `CREATE INDEX ON readings USING btree(device_uuid)`)
	if err != nil {
		return fmt.Errorf("failed to ensure device index for readings table: %w", err)
	}

	return nil
}

//InsertReading inserts a new reading
func (a *DBAgent) InsertReading(ctx context.Context, deviceUUID string, tCelcius json.Number, rh json.Number) (string, error) {
	row := a.db.QueryRowContext(ctx, `
INSERT INTO "readings" (
	"created_at",
	"modified_at",

	"device_uuid",
	"temperature_celcius",
	"relative_humidity"
) VALUES (
	NOW(), NOW(),
	$1, $2, $3
) RETURNING "uuid"`,
		deviceUUID, tCelcius, rh,
	)

	var uuid string
	err := row.Scan(&uuid)
	if err != nil {
		return "", fmt.Errorf("failed to insert into readings table: %w", err)
	}
	return uuid, nil
}
