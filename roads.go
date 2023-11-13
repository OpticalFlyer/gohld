package main

import (
	"net/http"

	"github.com/OpticalFlyer/hld/centerlines"
	"github.com/gin-gonic/gin"
)

func handleGetRoadCenterlineData(c *gin.Context) {
	var coords struct {
		South float64 `json:"south"`
		West  float64 `json:"west"`
		North float64 `json:"north"`
		East  float64 `json:"east"`
	}
	if err := c.ShouldBindJSON(&coords); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := centerlines.GetRoadCenterlineGeoJSON(coords.South, coords.West, coords.North, coords.East)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve road data"})
		return
	}

	c.Data(http.StatusOK, "application/json", data)
}
