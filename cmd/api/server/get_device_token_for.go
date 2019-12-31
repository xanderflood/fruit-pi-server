package server

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

//GetDeviceTokenRequest encodes a single request for user registration
type GetDeviceTokenRequest struct {
	DeviceUUID string `json:"device_uuid" binding:"required"`
}

//GetDeviceTokenFor generates a device request token for the specified device
func (a ServerAgent) GetDeviceTokenFor(c *gin.Context) {
	authorization, ok := a.authorize(c)
	if !ok {
		return
	}

	if !authorization.Admin {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "this endpoint is only accessible to users with administrative priveleges"})
		return
	}

	var req GetDeviceTokenRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device, err := a.dbClient.GetDeviceByUUID(c, req.DeviceUUID)
	if err != nil {
		err = fmt.Errorf("get device failed: %w", err)
		a.logger.Errorf(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub.dvc":   req.DeviceUUID,
		"frtpi.dvc": true,
	}).SignedString([]byte(a.jwtSigningSecret))
	if err != nil {
		err = fmt.Errorf("failed generating token: %w", err)
		a.logger.Errorf(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_uuid": req.DeviceUUID,
		"name":        device.Name,
		"token":       tokenString,
	})
}
