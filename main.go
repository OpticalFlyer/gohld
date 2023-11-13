package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/OpticalFlyer/hld/centerlines"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

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
	// Establish your database connection
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")
	var sslmode string

	if host == "localhost" {
		sslmode = "disable"
	} else {
		sslmode = "require"
	}

	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s password=%s", host, port, user, dbname, sslmode, password)

	var err error
	db, err = gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Failed to connect to database")
		return
	}
	defer db.Close()

	r := gin.Default()

	store := cookie.NewStore([]byte("malloc32$"))
	r.Use(sessions.Sessions("hld-session", store))

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.GET("/check-auth", func(c *gin.Context) {
		session := sessions.Default(c)
		authenticated := session.Get("authenticated")
		if authenticated == true {
			c.JSON(http.StatusOK, gin.H{"authenticated": true})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false})
		}
	})

	r.POST("/login", func(c *gin.Context) {
		var loginForm struct {
			Username string `form:"username"`
			Password string `form:"password"`
		}

		// This will bind the form data to the struct
		if err := c.ShouldBind(&loginForm); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check the credentials: if they are valid, set the session
		if loginForm.Username == "test" && loginForm.Password == "dummy" {
			session := sessions.Default(c)
			session.Set("authenticated", true)
			if err := session.Save(); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		}
	})

	r.POST("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete("authenticated")
		session.Save()
		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	})

	// Serve all files in the "static" directory
	r.Static("/static", "./static")

	r.POST("/road-centerline", handleGetRoadCenterlineData)

	r.POST("/sites-in-polygon", handleGetSitesInPolygon)

	r.Run() // listen and serve on 0.0.0.0:8080
}
