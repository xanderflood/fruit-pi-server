package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/xanderflood/fruit-pi-server/lib/api"
)

//InsertReading records a new reading from this device
func (a ServerAgent) InsertReading(c *gin.Context) {
	authorization, ok := a.authorize(c)
	if !ok {
		return
	}

	var req api.InsertReadingRequest
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

	c.JSON(http.StatusOK, api.Reading{
		DeviceUUID:  authorization.DeviceUUID,
		ReadingUUID: uuid,
	})
}
