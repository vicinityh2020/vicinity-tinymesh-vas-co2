package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type EventData struct {
	Value     int    `json:"value"`
	Unit      string `json:"unit"`
	Timestamp string `json:"timestamp"`
}

func (server *Server) vicinityEventHandler(c *gin.Context) {
	var e EventData

	oid, exists := c.Params.Get("oid")
	if !exists {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := c.BindJSON(&e); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// TODO: store as ONE of many-valueReadings-to-one-oid entry in database
	fmt.Println(oid)
	c.JSON(http.StatusOK, nil)
}

func (server *Server) handleTD(c *gin.Context) {
	c.JSON(http.StatusOK, server.vicinity.GetThingDescription())
}
