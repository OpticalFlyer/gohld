package main

import (
	"fmt"
	"log"
	"os"

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

func main() {
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

	db.AutoMigrate(&User{})

	r := gin.Default()

	store := cookie.NewStore([]byte("malloc32$"))
	r.Use(sessions.Sessions("hld-session", store))

	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.GET("/check-auth", checkAuth)
	r.POST("/login", login)
	r.POST("/logout", logout)

	r.POST("/road-centerline", handleGetRoadCenterlineData)
	r.POST("/sites-in-polygon", handleGetSitesInPolygon)

	// Serve all files in the "static" directory
	r.Static("/static", "./static")

	r.POST("/createuser", createUser)
	r.GET("/users", getUsers)
	r.DELETE("/users/:id", deleteUser)

	r.Run() // listen and serve on 0.0.0.0:8080
}
