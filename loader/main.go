package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Site struct {
	ID          uint   `gorm:"primary_key"`
	SiteID      string `gorm:"column:site_id;type:varchar(100)"`
	Address     string `gorm:"type:varchar(500)"`
	Eligibility string `gorm:"type:varchar(500)"`
	Tier        string `gorm:"type:varchar(500)"`
	Latitude    float64
	Longitude   float64
	AddressType string `gorm:"column:address_type;type:varchar(500)"`
	Geom        string `gorm:"type:geometry(Point, 4326)"`
}

func (Site) TableName() string {
	return "sites"
}

func main() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	dbname := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	//connectionString := "host=localhost user=hld dbname=hld sslmode=disable password=malloc32$" // Local
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=require password=%s", host, port, user, dbname, password)

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// AutoMigrate will create the table and add missing columns, indexes
	db.AutoMigrate(&Site{})

	// Add a GiST index for the Geom column
	db.Exec("CREATE INDEX idx_sites_geom ON sites USING GIST (geom);")

	// Parse the CSV and insert into the database
	parseCSVAndInsert(db)
}

func parseCSVAndInsert(db *gorm.DB) {
	f, err := os.Open("All Colorado Address Points CPF.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	totalLines := 0
	for scanner.Scan() {
		totalLines++
	}
	f.Seek(0, io.SeekStart) // Reset file reader to the beginning
	fmt.Println("Total lines in file:", totalLines)

	r := csv.NewReader(f)

	lineCount := 1
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading record on line %d: %v", lineCount, err)
			lineCount++
			continue
		}

		lat, err := strconv.ParseFloat(record[4], 64)
		if err != nil {
			log.Printf("Error parsing latitude on line %d: %v", lineCount, err)
			lineCount++
			continue
		}

		lon, err := strconv.ParseFloat(record[5], 64)
		if err != nil {
			log.Printf("Error parsing longitude on line %d: %v", lineCount, err)
			lineCount++
			continue
		}

		site := Site{
			SiteID:      record[0],
			Address:     record[1],
			Eligibility: record[2],
			Tier:        record[3],
			Latitude:    lat,
			Longitude:   lon,
			AddressType: record[6],
			Geom:        fmt.Sprintf("SRID=4326;POINT(%f %f)", lon, lat),
		}
		db.Create(&site)

		lineCount++
		progress := float64(lineCount) / float64(totalLines) * 100.0
		fmt.Printf("\rProgress: %.2f%%", progress)
	}
}
