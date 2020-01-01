package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xanderflood/fruit-pi-server/lib/api"
)

//RegisterDevice registers a device and responds with the newly generated UUID
func (a ServerAgent) RegisterDevice(c *gin.Context) {
	authorization, ok := a.authorize(c)
	if !ok {
		return
	}

	if !authorization.Admin {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "this endpoint is only accessible to users with administrative priveleges"})
		return
	}

	var req api.RegistrationRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uuid, err := a.dbClient.RegisterDevice(c, req.Name, string(req.Config))
	if err != nil {
		err = fmt.Errorf("register device failed: %w", err)
		a.logger.Errorf(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.Device{
		DeviceUUID: uuid,
		Name:       &req.Name,
		Config:     &req.Config,
	})
}
