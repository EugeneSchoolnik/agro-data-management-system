package repository

import (
	"agro-data-management-system/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSensorRepository_Lifecycle(t *testing.T) {
	// Порядок очищення важливий через FOREIGN KEY (спочатку дочірні, потім батьківські)
	db := setupTestDB(t, "sensors", "fields", "crops")
	defer db.Close()

	cropRepo := NewCropPostgres(db)
	fieldRepo := NewFieldPostgres(db)
	sensorRepo := NewSensorPostgres(db)

	// 1. Готуємо інфраструктуру: Культура -> Поле
	cropID, _ := cropRepo.Create(models.Crop{Name: "Пшениця"})
	fieldID, _ := fieldRepo.Create(models.Field{Name: "Ділянка А", CropID: cropID})

	// 2. ТЕСТ: Створення датчика
	newSensor := models.Sensor{
		FieldID:    fieldID,
		SensorType: "humidity",
		Status:     "active",
	}

	sensorID, err := sensorRepo.Create(newSensor)
	assert.NoError(t, err)
	assert.Greater(t, sensorID, 0)

	// 3. ТЕСТ: Отримання списку датчиків для поля
	sensors, err := sensorRepo.GetByFieldID(fieldID)
	assert.NoError(t, err)
	assert.Len(t, sensors, 1)
	assert.Equal(t, "humidity", sensors[0].SensorType)

	// 4. ТЕСТ: Оновлення статусу (імітація виходу з ладу або синхронізації)
	err = sensorRepo.UpdateStatus(sensorID, "inactive")
	assert.NoError(t, err)

	updated, _ := sensorRepo.GetByID(sensorID)
	assert.Equal(t, models.StatusInactive, updated.Status)
	assert.NotNil(t, updated.LastSync)

	// 5. ТЕСТ: Видалення
	err = sensorRepo.Delete(sensorID)
	assert.NoError(t, err)
}
