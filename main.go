package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GeoJsonStruct struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Crs  struct {
		Type       string `json:"type"`
		Properties struct {
			Name string `json:"name"`
		} `json:"properties"`
	} `json:"crs"`
	Features []struct {
		Type       string `json:"type"`
		Properties struct {
			ID              int         `json:"_id"`
			PARCELID        int         `json:"PARCELID"`
			FEATURETYPE     string      `json:"FEATURE_TYPE"`
			DATEEFFECTIVE   string      `json:"DATE_EFFECTIVE"`
			DATEEXPIRY      string      `json:"DATE_EXPIRY"`
			PLANID          int         `json:"PLANID"`
			PLANDESCRIPTION string      `json:"PLAN_DESCRIPTION"`
			PLANNAME        string      `json:"PLAN_NAME"`
			PLANTYPE        string      `json:"PLAN_TYPE"`
			STATEDAREA      string      `json:"STATEDAREA"`
			SOURCEID        interface{} `json:"SOURCE_ID"`
			ADDRESSPOINTID  int         `json:"ADDRESS_POINT_ID"`
			ADDRESSNUMBER   string      `json:"ADDRESS_NUMBER"`
			LINEARNAMEID    int         `json:"LINEAR_NAME_ID"`
			LINEARNAMEFULL  string      `json:"LINEAR_NAME_FULL"`
			AROLLSOURCEDESC string      `json:"AROLL_SOURCE_DESC"`
			ADDRESSID       int         `json:"ADDRESS_ID"`
			OBJECTID        string      `json:"OBJECTID"`
			TRANSIDCREATE   float64     `json:"TRANS_ID_CREATE"`
			TRANSIDEXPIRE   float64     `json:"TRANS_ID_EXPIRE"`
		} `json:"properties"`
		Geometry struct {
			Type        string          `json:"type"`
			Coordinates [][][][]float64 `json:"coordinates"`
		} `json:"geometry"`
	} `json:"features"`
}

func main() {

	torontoContent, err := os.ReadFile("torontoeverything.json")
	if err != nil {
		fmt.Println("Error reading torontoeverything.json:", err)
	}

	var jsonData GeoJsonStruct
	err = json.Unmarshal(torontoContent, &jsonData)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
	}

	for i := 0; i < 122; i++ {
		// create a new file
		savedContent, err := os.Create(fmt.Sprintf("json/toronto_%d.json", i))
		if err != nil {
			fmt.Println("Error creating file:", err)
		}

		defer savedContent.Close()
	}

	totalFeatures := len(jsonData.Features)
	featuresPerFile := 4500

	fmt.Printf("Total features: %d\n", totalFeatures)
	fmt.Printf("Features per file: %d\n", featuresPerFile)

	for fileIndex := 0; fileIndex < 122; fileIndex++ {
		startIndex := fileIndex * featuresPerFile
		endIndex := startIndex + featuresPerFile

		// Make sure we don't go beyond the total features
		if startIndex >= totalFeatures {
			break
		}
		if endIndex > totalFeatures {
			endIndex = totalFeatures
		}

		// Create new structure with same type, name, and crs
		splitData := GeoJsonStruct{
			Type:     jsonData.Type,
			Name:     jsonData.Name,
			Crs:      jsonData.Crs,
			Features: jsonData.Features[startIndex:endIndex],
		}

		// Marshal to JSON
		jsonBytes, err := json.MarshalIndent(splitData, "", "  ")
		if err != nil {
			fmt.Printf("Error marshalling data for file %d: %v\n", fileIndex, err)
			continue
		}

		// Write to file
		filename := fmt.Sprintf("json/toronto_%d.json", fileIndex)
		err = os.WriteFile(filename, jsonBytes, 0644)
		if err != nil {
			fmt.Printf("Error writing file %s: %v\n", filename, err)
			continue
		}

		fmt.Printf("Created %s with %d features\n", filename, len(splitData.Features))
	}

	fmt.Println("File splitting completed!")

}
