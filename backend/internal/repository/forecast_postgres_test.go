package repository

import (
	"testing"
	"time"

	"agro-data-management-system/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestForecastRepository_Lifecycle(t *testing.T) {
	// Очищуємо все дерево залежностей
	db := setupTestDB(t, "forecasts", "pests", "fields", "crops")
	defer db.Close()

	cropRepo := NewCropPostgres(db)
	fieldRepo := NewFieldPostgres(db)
	pestRepo := NewPestPostgres(db)
	forecastRepo := NewForecastPostgres(db)

	// 1. Створюємо оточення
	cID, _ := cropRepo.Create(models.Crop{Name: "Пшениця"})
	fID, _ := fieldRepo.Create(models.Field{Name: "Сектор 7", CropID: cID})
	pID, _ := pestRepo.Create(models.Pest{Name: "Черепашка", ScientificName: "E. integriceps"})

	// 2. ТЕСТ: Створення прогнозу
	now := time.Now().UTC().Truncate(time.Second)
	newForecast := models.Forecast{
		FieldID:        fID,
		PestID:         pID,
		Probability:    0.85,
		Recommendation: "Необхідна термінова обробка інсектицидами",
		CreatedAt:      now,
	}

	id, err := forecastRepo.Create(newForecast)
	assert.NoError(t, err)
	assert.Greater(t, id, 0)

	// 3. ТЕСТ: Отримання останнього прогнозу
	latest, err := forecastRepo.GetLatestByField(fID)
	assert.NoError(t, err)
	assert.Equal(t, 0.85, latest.Probability)
	assert.Equal(t, "Необхідна термінова обробка інсектицидами", latest.Recommendation)

	// 4. ТЕСТ: Історія прогнозів
	history, err := forecastRepo.GetHistoryByField(fID)
	assert.NoError(t, err)
	assert.Len(t, history, 1)
}
