package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xanderflood/fruit-pi-server/lib/api"
)

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

	var req api.ConfigureDeviceRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uuid, err := a.dbClient.ConfigureDevice(c, req.DeviceUUID, req.Name, string(req.Config))
	if err != nil {
		err = fmt.Errorf("configure device failed: %w", err)
		a.logger.Errorf(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.Device{
		DeviceUUID: req.DeviceUUID,
		Name:       &uuid,
		Config:     &req.Config,
	})
}
