package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//RegistrationRequest encodes a single request for user registration
type RegistrationRequest struct {
	Name   string          `json:"name" binding:"required"`
	Config json.RawMessage `json:"config" binding:"required"`
}

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

	var req RegistrationRequest
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

	c.JSON(http.StatusOK, gin.H{
		"device_uuid": uuid,
		"name":        req.Name,
		"config":      req.Config,
	})
}
