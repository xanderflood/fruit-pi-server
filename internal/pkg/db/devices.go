package db

import (
	"context"
	"fmt"
)

//RegisterDevice registers a new device in the table
func (a *DBAgent) RegisterDevice(ctx context.Context, name string, config string) (string, error) {
	row := a.db.QueryRowContext(ctx, `
INSERT INTO "devices" (
	"created_at",
	"modified_at",

	"name",
	"config"
) VALUES (
	NOW(), NOW(),
	$1, $2
) RETURNING "uuid"`,
		name, config,
	)

	var uuid string
	err := row.Scan(&uuid)
	if err != nil {
		return "", fmt.Errorf("failed to insert into devices table: %w", err)
	}
	return uuid, nil
}

//EnsureDevicesTable EnsureDevicesTable
func (a *DBAgent) EnsureDevicesTable(ctx context.Context) error {
	_, err := a.db.ExecContext(ctx, `
CREATE TABLE IF NOT EXISTS "devices"
(	"uuid" UUID DEFAULT gen_random_uuid(),
	"created_at" timestamp NOT NULL,
	"modified_at" timestamp NOT NULL,
	"deleted_at" timestamp,

	"name" varchar NOT NULL,
	"config" text NOT NULL,
	PRIMARY KEY ("uuid")
)`)
	if err != nil {
		return fmt.Errorf("failed to ensure accounts table: %w", err)
	}
	return nil
}

//ConfigureDevice updates a device's configuration
func (a *DBAgent) ConfigureDevice(ctx context.Context, uuid string, name string, config string) (string, error) {
	_, err := a.db.ExecContext(ctx, `
INSERT INTO "devices" (
	"created_at",
	"modified_at",

	"uuid",
	"name",
	"config"
) VALUES (
	NOW(), NOW(),
	$1, $2, $3
)

ON CONFLICT (uuid) DO UPDATE
SET "name" = $2,
	"config" = $3
RETURNING "uuid"`,
		uuid, name, config,
	)
	if err != nil {
		return "", fmt.Errorf("failed to update device config: %w", err)
	}
	return uuid, nil
}

//GetDeviceByUUID gets a device by UUID
func (a *DBAgent) GetDeviceByUUID(ctx context.Context, uuid string) (device Device, err error) {
	row := a.db.QueryRowContext(ctx, fmt.Sprintf(`
SELECT %s FROM "devices"
WHERE
	"deleted_at" IS NULL
	AND
	"uuid" = $1
`, StandardDeviceFieldNameList),
		uuid,
	)

	err = row.Scan(device.StandardFieldPointers()...)
	if err != nil {
		err = fmt.Errorf("failed to get device by uuid `%s`: %w", uuid, err)
	}
	return
}
