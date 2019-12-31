package db

import (
	"encoding/json"
	"time"
)

//Model contains generic fields shared by all models
type Model struct {
	UUID       string     `json:"uuid"`
	CreatedAt  time.Time  `json:"created_at"`
	ModifiedAt time.Time  `json:"modified_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

//Device represents a single fruiting chamber device
type Device struct {
	Model

	Name   string `json:"name"`
	Config string `json:"config"`
}

const StandardDeviceFieldNameList = `
	"uuid",
	"created_at",
	"modified_at",

	"name",
	"config"
`

func (d *Device) StandardFieldPointers() []interface{} {
	return []interface{}{
		&d.UUID,
		&d.CreatedAt,
		&d.ModifiedAt,

		&d.Name,
		&d.Config,
	}
}

//Reading represents a single sensor reading from the chamber
type Reading struct {
	Model

	DeviceUUID         string      `json:"device_uuid"`
	TemperatureCelcius json.Number `json:"temperature_celcius"`
	RelativeHumidity   json.Number `json:"relative_humidity"`
}

const StandardReadingFieldNameList = `
	"uuid",
	"created_at",
	"modified_at",

	"device_uuid",
	"temperature_celcius",
	"relative_humidity",
`

func (r *Reading) StandardFieldPointers() []interface{} {
	return []interface{}{
		&r.UUID,
		&r.CreatedAt,
		&r.ModifiedAt,

		&r.DeviceUUID,
		&r.TemperatureCelcius,
		&r.RelativeHumidity,
	}
}
