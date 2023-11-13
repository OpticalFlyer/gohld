package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleGetSitesInPolygon(c *gin.Context) {
	var polyRequest PolygonRequest
	if err := c.ShouldBindJSON(&polyRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert the provided polygon coordinates to WKT (Well-Known Text) format
	polygonWKT := "POLYGON(("
	for _, coord := range polyRequest.Polygon {
		polygonWKT += fmt.Sprintf("%f %f,", coord[1], coord[0]) // Note: PostGIS uses Longitude Latitude ordering
	}
	// Close the polygon by adding the first point again
	polygonWKT += fmt.Sprintf("%f %f", polyRequest.Polygon[0][1], polyRequest.Polygon[0][0])
	polygonWKT += "))"

	fmt.Println(polygonWKT)

	// Fetch points within the polygon and convert to GeoJSON
	rows, err := db.Raw(`SELECT ST_AsGeoJSON(geom) FROM sites WHERE ST_Within(geom, ST_GeomFromText(?, 4326))`, polygonWKT).Rows()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch data from database"})
		return
	}
	defer rows.Close()

	fmt.Println("Fetched rows")

	// Convert rows to a GeoJSON feature collection
	featureCollection := map[string]interface{}{
		"type":     "FeatureCollection",
		"features": []map[string]interface{}{},
	}

	for rows.Next() {
		var geojsonStr string
		rows.Scan(&geojsonStr)
		var feature map[string]interface{}
		json.Unmarshal([]byte(geojsonStr), &feature)
		featureCollection["features"] = append(featureCollection["features"].([]map[string]interface{}), feature)
	}

	c.JSON(http.StatusOK, featureCollection)
}
