package repository

import (
	"testing"
	"time"

	"agro-data-management-system/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestMetricRepository_Lifecycle(t *testing.T) {
	// Очищуємо все дерево залежностей
	db := setupTestDB(t, "metrics", "sensors", "fields", "crops")
	defer db.Close()

	cropRepo := NewCropPostgres(db)
	fieldRepo := NewFieldPostgres(db)
	sensorRepo := NewSensorPostgres(db)
	metricRepo := NewMetricPostgres(db)

	// 1. Створюємо оточення
	cID, _ := cropRepo.Create(models.Crop{Name: "Пшениця"})
	fID, _ := fieldRepo.Create(models.Field{Name: "Поле 1", CropID: cID})
	sID, _ := sensorRepo.Create(models.Sensor{FieldID: fID, SensorType: "DHT22", Status: models.StatusActive})

	// 2. ТЕСТ: Створення метрики
	now := time.Now().UTC().Truncate(time.Second)
	m := models.Metric{
		SensorID:   sID,
		Value:      26.7,
		RecordedAt: now,
	}

	id, err := metricRepo.Create(m)
	assert.NoError(t, err)
	assert.Greater(t, id, int64(0))

	// 3. ТЕСТ: Отримання останньої
	latest, err := metricRepo.GetLatestBySensor(sID)
	assert.NoError(t, err)
	assert.Equal(t, 26.7, latest.Value)
	assert.WithinDuration(t, now, latest.RecordedAt.UTC(), time.Second)

	// 4. ТЕСТ: Історія
	history, err := metricRepo.GetHistoryBySensor(sID, now.Add(-1*time.Minute), now.Add(1*time.Minute))
	assert.NoError(t, err)
	assert.Len(t, history, 1)
}
