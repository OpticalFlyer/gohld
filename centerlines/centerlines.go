package centerlines

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

// GeoJSON model
type GeoJSON struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string     `json:"type"`
	Geometry   Geometry   `json:"geometry"`
	Properties Properties `json:"properties"`
}

type Geometry struct {
	Type        string      `json:"type"`
	Coordinates [][]float64 `json:"coordinates"`
}

type Properties struct {
	ID   int64             `json:"id"`
	Tags map[string]string `json:"tags"`
}

// OSMData is the JSON structure returned by the Overpass API
type Element struct {
	Type  string            `json:"type"`
	ID    int64             `json:"id"`
	Nodes []int64           `json:"nodes,omitempty"`
	Tags  map[string]string `json:"tags,omitempty"`
	Lat   float64           `json:"lat,omitempty"`
	Lon   float64           `json:"lon,omitempty"`
}

type OSMData struct {
	Elements []Element `json:"elements"`
}

const overpassAPIURL = "https://overpass-api.de/api/interpreter"

func GetRoadCenterlineGeoJSON(south, west, north, east float64) ([]byte, error) {
	body, err := getRoadCenterlineData(south, west, north, east)
	if err != nil {
		return nil, err
	}

	// Unmarshal into OSMData struct
	var osmData OSMData
	err = json.Unmarshal(body, &osmData)
	if err != nil {
		return nil, err
	}

	// Map node IDs to their lat/lon for easy lookup
	nodeMap := make(map[int64][2]float64)
	for _, elem := range osmData.Elements {
		if elem.Type == "node" {
			nodeMap[elem.ID] = [2]float64{elem.Lon, elem.Lat}
		}
	}

	// Build GeoJSON
	geoJSON := &GeoJSON{
		Type:     "FeatureCollection",
		Features: []Feature{},
	}

	for _, elem := range osmData.Elements {
		if elem.Type == "way" {
			coords := make([][]float64, len(elem.Nodes))
			for i, nodeID := range elem.Nodes {
				if coord, ok := nodeMap[nodeID]; ok {
					coords[i] = coord[:]
				}
			}

			feature := Feature{
				Type: "Feature",
				Geometry: Geometry{
					Type:        "LineString",
					Coordinates: coords,
				},
				Properties: Properties{
					ID:   elem.ID,
					Tags: elem.Tags,
				},
			}

			geoJSON.Features = append(geoJSON.Features, feature)
		}
	}

	// Serialize the data to JSON
	jsonData, err := json.MarshalIndent(geoJSON, "", "  ")
	if err != nil {
		log.Fatalf("Error serializing to JSON: %v", err)
	}

	return jsonData, nil
}

func getRoadCenterlineData(south, west, north, east float64) ([]byte, error) {
	query := fmt.Sprintf(`
		[out:json][timeout:25];
		(
			way["highway"]["highway"!~"proposed"]["highway"!~"service"]["highway"!~"path"](%f,%f,%f,%f);
		);
		out body;
		>;
		out skel qt;
	`, south, west, north, east)

	resp, err := http.PostForm(overpassAPIURL, url.Values{"data": {query}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
