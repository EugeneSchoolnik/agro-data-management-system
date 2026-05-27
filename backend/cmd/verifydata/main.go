package main

import (
	"agro-data-management-system/internal/config"
	"agro-data-management-system/internal/repository"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig("config/local.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := repository.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	repos := repository.NewRepositories(db)

	// Get Field 2 details
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("FIELD DETAILS (ID: 2)")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fieldWithCrop, err := repos.Field.GetByIDWithCrop(2)
	if err != nil {
		log.Printf("Error getting field: %v", err)
	} else {
		fmt.Printf("Field Name: %s\n", fieldWithCrop.Name)
		fmt.Printf("Area: %.2f ha\n", fieldWithCrop.Area)
		fmt.Printf("Location: %s\n", fieldWithCrop.Location)
		fmt.Printf("Crop: %s (%s)\n", fieldWithCrop.CropName, fieldWithCrop.CropVariety)
		fmt.Printf("Created At: %s\n", fieldWithCrop.CreatedAt)
	}

	// Get sensors for Field 2
	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("SENSORS (Field ID: 2)")
	fmt.Println("═══════════════════════════════════════════════════════════")
	sensors, err := repos.Sensor.GetByFieldID(2)
	if err != nil {
		log.Printf("Error getting sensors: %v", err)
	} else {
		fmt.Printf("Total Sensors: %d\n", len(sensors))
		for _, s := range sensors {
			fmt.Printf("  ID: %d | Type: %s | Status: %s | Last Sync: %s\n",
				s.ID, s.SensorType, s.Status, s.LastSync)
		}
	}

	// Get metrics for each sensor
	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("METRICS")
	fmt.Println("═══════════════════════════════════════════════════════════")
	for _, s := range sensors {
		// Get latest metric
		latestMetric, err := repos.Metric.GetLatestBySensor(s.ID)
		if err != nil {
			continue
		}
		fmt.Printf("Sensor %d (%s) - Latest Value: %.2f @ %s\n",
			s.ID, s.SensorType, latestMetric.Value, latestMetric.RecordedAt)
	}

	// Get pests
	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("PESTS")
	fmt.Println("═══════════════════════════════════════════════════════════")
	pests, err := repos.Pest.GetAll()
	if err != nil {
		log.Printf("Error getting pests: %v", err)
	} else {
		fmt.Printf("Total Pests: %d\n", len(pests))
		for _, p := range pests {
			fmt.Printf("  ID: %d | %s (%s)\n", p.ID, p.Name, p.ScientificName)
		}
	}

	// Get forecasts for Field 2
	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("FORECASTS (Field ID: 2)")
	fmt.Println("═══════════════════════════════════════════════════════════")
	history, err := repos.Forecast.GetHistoryByField(2)
	if err != nil {
		log.Printf("Error getting forecasts: %v", err)
	} else {
		fmt.Printf("Total Forecasts: %d\n", len(history))
		for _, f := range history {
			fmt.Printf("  ID: %d | Pest: %d | Probability: %.2f | Created: %s\n",
				f.ID, f.PestID, f.Probability, f.CreatedAt)
		}
	}

	// Count total records
	fmt.Println("\n═══════════════════════════════════════════════════════════")
	fmt.Println("DATABASE SUMMARY")
	fmt.Println("═══════════════════════════════════════════════════════════")

	var cropCount, fieldCount, sensorCount, metricCount, pestCount, forecastCount int64

	db.Get(&cropCount, "SELECT COUNT(*) FROM crops")
	db.Get(&fieldCount, "SELECT COUNT(*) FROM fields")
	db.Get(&sensorCount, "SELECT COUNT(*) FROM sensors")
	db.Get(&metricCount, "SELECT COUNT(*) FROM metrics")
	db.Get(&pestCount, "SELECT COUNT(*) FROM pests")
	db.Get(&forecastCount, "SELECT COUNT(*) FROM forecasts")

	fmt.Printf("Total Crops: %d\n", cropCount)
	fmt.Printf("Total Fields: %d\n", fieldCount)
	fmt.Printf("Total Sensors: %d\n", sensorCount)
	fmt.Printf("Total Metrics: %d\n", metricCount)
	fmt.Printf("Total Pests: %d\n", pestCount)
	fmt.Printf("Total Forecasts: %d\n", forecastCount)
}
