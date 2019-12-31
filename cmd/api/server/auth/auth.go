package auth

import (
	"errors"
	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/xanderflood/fruit-pi-server/pkg/db"
)

//Authorization describes the authorities stored in a user JWT
type Authorization struct {
	jwt.StandardClaims
	DeviceUUID string `json:"sub.dvc,omitempty"`
	Admin      bool   `json:"login.adm,omitempty"`
	Device     bool   `json:"frtpi.dvc,omitempty"`
}

//Valid applies standard JWT validations as well as generic
//user authorization rules.
func (a *Authorization) Valid() error {
	if err := a.StandardClaims.Valid(); err != nil {
		return err
	}

	if !a.Admin && len(a.DeviceUUID) == 0 {
		return errors.New("non-admin user does not have an identity")
	}

	return nil
}

func (a Authorization) GetDBAuthorization(deviceUUID string) (db.Authorization, error) {
	//For this API, devie UUID in the request body is always optional.
	//If it's missing, it'll be inferred from the token.
	if len(deviceUUID) == 0 {
		deviceUUID = a.DeviceUUID
	}

	if a.Admin || (a.Device && a.DeviceUUID == deviceUUID) {
		return db.Authorization{DeviceUUID: deviceUUID}, nil
	}
	return db.Authorization{}, fmt.Errorf("refusing to produce database authorization for device UUID `%s`", deviceUUID)
}
