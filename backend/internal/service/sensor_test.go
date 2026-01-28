package service

import (
	"errors"
	"testing"

	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupSensorServiceTest() (*sensorService, *mocks.SensorRepository, *mocks.FieldRepository) {
	sRepo := new(mocks.SensorRepository)
	fRepo := new(mocks.FieldRepository)
	log := zap.NewNop()
	srv := NewSensorService(sRepo, fRepo, log).(*sensorService)
	return srv, sRepo, fRepo
}

func TestSensorService_Register(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv, sRepo, fRepo := setupSensorServiceTest()
		sensor := models.Sensor{FieldID: 1, SensorType: "DHT22", Status: models.StatusActive}

		// Поле існує
		fRepo.On("GetByIDWithCrop", 1).Return(models.FieldWithCrop{}, nil).Once()
		// Створення успішне
		sRepo.On("Create", sensor).Return(models.Sensor{ID: 101, FieldID: 1, SensorType: "DHT22", Status: models.StatusActive}, nil).Once()

		result, err := srv.Register(sensor)

		assert.NoError(t, err)
		assert.Equal(t, models.Sensor{ID: 101, FieldID: 1, SensorType: "DHT22", Status: models.StatusActive}, result)
		mock.AssertExpectationsForObjects(t, sRepo, fRepo)
	})

	t.Run("Field_Not_Found", func(t *testing.T) {
		srv, sRepo, fRepo := setupSensorServiceTest()
		sensor := models.Sensor{FieldID: 404, SensorType: "Camera", Status: models.StatusActive}

		fRepo.On("GetByIDWithCrop", 404).Return(models.FieldWithCrop{}, errors.New("sql: no rows")).Once()

		result, err := srv.Register(sensor)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "field 404 not found")
		assert.Equal(t, sensor, result)
		sRepo.AssertNotCalled(t, "Create", mock.Anything)
	})

	t.Run("Invalid_Status_Enum", func(t *testing.T) {
		srv, sRepo, _ := setupSensorServiceTest()
		// Передаємо невалідний статус
		sensor := models.Sensor{FieldID: 1, SensorType: "DHT22", Status: models.SensorStatus("broken_string")}

		_, err := srv.Register(sensor)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid sensor status")
		sRepo.AssertNotCalled(t, "Create", mock.Anything)
	})

	t.Run("Validation_Error_Short_Type", func(t *testing.T) {
		srv, _, _ := setupSensorServiceTest()
		sensor := models.Sensor{FieldID: 1, SensorType: "it", Status: models.StatusActive} // "it" занадто коротко

		_, err := srv.Register(sensor)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "SensorType")
	})
}

func TestSensorService_UpdateStatus(t *testing.T) {
	t.Run("Valid_Transition", func(t *testing.T) {
		srv, sRepo, _ := setupSensorServiceTest()
		sRepo.On("UpdateStatus", 1, models.StatusError).Return(nil).Once()

		err := srv.UpdateStatus(1, models.StatusError)

		assert.NoError(t, err)
	})

	t.Run("Invalid_Status_Update", func(t *testing.T) {
		srv, sRepo, _ := setupSensorServiceTest()
		err := srv.UpdateStatus(1, models.SensorStatus("unknown"))

		assert.Error(t, err)
		sRepo.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything)
	})
}
