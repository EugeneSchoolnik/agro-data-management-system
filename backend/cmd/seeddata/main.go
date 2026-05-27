package main

import (
	"agro-data-management-system/internal/config"
	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config/local.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := repository.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	repos := repository.NewRepositories(db)

	// 0. Clean up existing data
	fmt.Println("Cleaning up existing data...")
	cleanupQueries := []string{
		"DELETE FROM forecasts",
		"DELETE FROM metrics",
		"DELETE FROM sensors",
		"DELETE FROM fields",
		"DELETE FROM crops",
		"DELETE FROM pests",
	}
	for _, query := range cleanupQueries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("failed to cleanup: %v", err)
		}
	}
	fmt.Println("✓ Tables cleaned up\n")

	// 1. Create Crop
	fmt.Println("Creating crop...")
	crop := models.Crop{
		Name:        "Пшениця",
		Variety:     "Україна",
		Description: "Озима пшениця українського сорту",
	}
	createdCrop, err := repos.Crop.Create(crop)
	if err != nil {
		log.Fatalf("failed to create crop: %v", err)
	}
	fmt.Printf("✓ Crop created: %+v\n", createdCrop)

	// 2. Create Field
	fmt.Println("\nCreating field...")
	field := models.Field{
		Name:     "Поле № 1",
		Area:     25.5,
		Location: "Київська область, село Петрівка",
		CropID:   &createdCrop.ID,
	}
	createdField, err := repos.Field.Create(field)
	if err != nil {
		log.Fatalf("failed to create field: %v", err)
	}
	fmt.Printf("✓ Field created: %+v\n", createdField)

	// 3. Create Sensors
	fmt.Println("\nCreating sensors...")
	sensorTypes := []string{"temperature", "humidity", "soil_moisture"}
	var sensorIDs []int

	for _, sensorType := range sensorTypes {
		sensor := models.Sensor{
			FieldID:    createdField.ID,
			SensorType: sensorType,
			Status:     models.StatusActive,
		}
		createdSensor, err := repos.Sensor.Create(sensor)
		if err != nil {
			log.Fatalf("failed to create sensor: %v", err)
		}
		sensorIDs = append(sensorIDs, createdSensor.ID)
		fmt.Printf("✓ Sensor created: %s (ID: %d)\n", sensorType, createdSensor.ID)
	}

	// 4. Create Pests
	fmt.Println("\nCreating pests...")
	pests := []models.Pest{
		{
			Name:           "Хлібна жужелиця",
			ScientificName: "Zabrus tenebrioides",
		},
		{
			Name:           "Гісяна совка",
			ScientificName: "Agrotis segetum",
		},
		{
			Name:           "Fusarium",
			ScientificName: "Fusarium graminearum",
		},
	}

	var pestIDs []int
	for _, pest := range pests {
		createdPest, err := repos.Pest.Create(pest)
		if err != nil {
			log.Fatalf("failed to create pest: %v", err)
		}
		pestIDs = append(pestIDs, createdPest.ID)
		fmt.Printf("✓ Pest created: %s (ID: %d)\n", pest.Name, createdPest.ID)
	}

	// 5. Create Metrics for last 7 days
	fmt.Println("\nCreating metrics for last 7 days...")
	now := time.Now()
	sevenDaysAgo := now.AddDate(0, 0, -7)

	// Generate metrics every 3 hours for each sensor
	for i := 0; i < 56; i++ { // 56 * 3 hours = 7 days
		recordedAt := sevenDaysAgo.Add(time.Hour * 3 * time.Duration(i))

		for j, sensorID := range sensorIDs {
			var value float64
			// Generate realistic values based on sensor type
			switch sensorTypes[j] {
			case "temperature":
				value = 15 + rand.Float64()*10 // 15-25°C
			case "humidity":
				value = 40 + rand.Float64()*40 // 40-80%
			case "soil_moisture":
				value = 30 + rand.Float64()*30 // 30-60%
			}

			metric := models.Metric{
				SensorID:   sensorID,
				Value:      value,
				RecordedAt: recordedAt,
			}
			_, err := repos.Metric.Create(metric)
			if err != nil {
				log.Fatalf("failed to create metric: %v", err)
			}
		}
	}
	fmt.Printf("✓ Created metrics: 56 records per sensor (168 total records)\n")

	// 6. Create Forecasts
	fmt.Println("\nCreating forecasts...")
	for _, pestID := range pestIDs {
		forecast := models.Forecast{
			FieldID:        createdField.ID,
			PestID:         pestID,
			Probability:    0.3 + rand.Float64()*0.5, // 0.3-0.8
			Recommendation: "Рекомендується проведення обробки інсектицидом у найближчі 2-3 дні",
			CreatedAt:      now,
		}
		createdForecast, err := repos.Forecast.Create(forecast)
		if err != nil {
			log.Fatalf("failed to create forecast: %v", err)
		}
		fmt.Printf("✓ Forecast created (ID: %d): Pest ID %d, Probability %.2f\n", createdForecast.ID, pestID, createdForecast.Probability)
	}

	fmt.Println("\n✓ All data created successfully!")
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  - Crop ID: %d\n", createdCrop.ID)
	fmt.Printf("  - Field ID: %d\n", createdField.ID)
	fmt.Printf("  - Sensors: %d\n", len(sensorIDs))
	fmt.Printf("  - Pests: %d\n", len(pestIDs))
	fmt.Printf("  - Metrics: 168\n")
	fmt.Printf("  - Forecasts: %d\n", len(pestIDs))
}
