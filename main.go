package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/OpticalFlyer/hld/centerlines"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PolygonRequest struct {
	Polygon [][]float64 `json:"polygon"`
}

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

	// Establish your database connection
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	//connectionString := "host=localhost user=hld dbname=hld sslmode=disable password=malloc32$" // Local
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=require password=%s", host, port, user, dbname, password)

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}
	defer db.Close()
	fmt.Println("Connected to database")

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

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// Serve all files in the "static" directory
	r.Static("/static", "./static")

	r.POST("/road-centerline", handleGetRoadCenterlineData)

	r.POST("/sites-in-polygon", handleGetSitesInPolygon)

	r.Run() // listen and serve on 0.0.0.0:8080
}
