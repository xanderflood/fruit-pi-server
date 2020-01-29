package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xanderflood/fruit-pi-server/lib/api"
)

//
// TODO: rename to simply get-device and do a little refactoring
//  so that new fields are always included on all endpoints - that'll
//  help avoid spurious diffs in the terraform provider.
//

//GetDeviceConfig gets the current configuration text for the device
func (a ServerAgent) GetDeviceConfig(c *gin.Context) {
	authorization, ok := a.authorize(c)
	if !ok {
		return
	}

	var req api.GetDeviceConfigRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var uuid string
	if authorization.Admin && req.DeviceUUID != nil {
		uuid = *req.DeviceUUID
	} else {
		uuid = authorization.DeviceUUID
	}

	device, err := a.dbClient.GetDeviceByUUID(c, uuid)
	if err != nil {
		err = fmt.Errorf("get device failed: %w", err)
		a.logger.Errorf(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_uuid": uuid,
		"name":        device.Name,
		"config":      json.RawMessage(device.Config),
	})
}
