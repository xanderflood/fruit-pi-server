package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//GetDeviceConfig gets the current configuration text for the device
func (a ServerAgent) GetDeviceConfig(c *gin.Context) {
	authorization, ok := a.authorize(c)
	if !ok {
		return
	}

	//no body - the JWT needs to contain the device UUID

	device, err := a.dbClient.GetDeviceByUUID(c, authorization.DeviceUUID)
	if err != nil {
		err = fmt.Errorf("get device failed: %w", err)
		a.logger.Errorf(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_uuid": authorization.DeviceUUID,
		"config":      json.RawMessage(device.Config),
	})
}
