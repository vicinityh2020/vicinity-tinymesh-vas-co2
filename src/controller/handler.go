package controller

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"vicinity-tinymesh-vas-co2/src/model"
)

type EventData struct {
	Value     int    `json:"value" binding:"required"`
	Unit      string `json:"unit"`
	Timestamp string `json:"timestamp" binding:"required"`
}

func (server *Server) vicinityEventHandler(c *gin.Context) {
	oid, exists := c.Params.Get("oid")
	if !exists {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	eid, exists := c.Params.Get("eid")
	if !exists {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var sensor model.Sensor

	id, err := uuid.FromString(oid)

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	server.db.Where(model.Sensor{Oid: id}).FirstOrCreate(&sensor, model.Sensor{Oid: id, Eid: eid, Unit: "ppm"})

	var e EventData
	if err := c.BindJSON(&e); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	server.db.Create(&model.Reading{Value: e.Value, Timestamp: e.Timestamp, SensorOid: sensor.Oid})

	c.JSON(http.StatusOK, nil)
}

func (server *Server) handleTD(c *gin.Context) {
	c.JSON(http.StatusOK, server.vicinity.GetThingDescription())
}
