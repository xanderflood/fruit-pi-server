package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//ConfigureDeviceRequest encodes a single request for user registration
type ConfigureDeviceRequest struct {
	DeviceUUID string          `json:"device_uuid" binding:"required"`
	Config     json.RawMessage `json:"config" binding:"required"`
}

//ConfigureDevice overwrites the configuration text for the device
func (a ServerAgent) ConfigureDevice(c *gin.Context) {
	authorization, ok := a.authorize(c)
	if !ok {
		return
	}

	if !authorization.Admin {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "this endpoint is only accessible to users with administrative priveleges"})
		return
	}

	var req ConfigureDeviceRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = a.dbClient.ConfigureDevice(c, req.DeviceUUID, string(req.Config))
	if err != nil {
		err = fmt.Errorf("configure device failed: %w", err)
		a.logger.Errorf(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_uuid": req.DeviceUUID,
		"config":      req.Config,
	})
}
