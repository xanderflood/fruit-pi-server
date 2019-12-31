package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//InsertReadingRequest encodes a single request for user registration
type InsertReadingRequest struct {
	TemperatureCelcius json.Number `json:"temperature_celcius" binding:"required"`
	RelativeHumidity   json.Number `json:"relative_humidity" binding:"required"`
}

//InsertReading records a new reading from this device
func (a ServerAgent) InsertReading(c *gin.Context) {
	authorization, ok := a.authorize(c)
	if !ok {
		return
	}

	var req InsertReadingRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uuid, err := a.dbClient.InsertReading(c, authorization.DeviceUUID, req.TemperatureCelcius, req.RelativeHumidity)
	if err != nil {
		err = fmt.Errorf("configure device failed: %w", err)
		a.logger.Errorf(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device_uuid":  authorization.DeviceUUID,
		"reading_uuid": uuid,
	})
}
