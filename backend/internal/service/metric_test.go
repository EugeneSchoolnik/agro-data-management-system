package service

import (
	"errors"
	"testing"
	"time"

	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func setupMetricServiceTest() (*metricService, *mocks.MetricRepository, *mocks.SensorRepository) {
	mRepo := new(mocks.MetricRepository)
	sRepo := new(mocks.SensorRepository)
	log := zap.NewNop()
	srv := NewMetricService(mRepo, sRepo, log).(*metricService)
	return srv, mRepo, sRepo
}

func TestMetricService_Save(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv, mRepo, sRepo := setupMetricServiceTest()
		m := models.Metric{SensorID: 1, Value: 25.5}

		// Датчик існує і він Active
		sRepo.On("GetByID", 1).Return(models.Sensor{ID: 1, Status: models.StatusActive}, nil).Once()
		mRepo.On("Create", m).Return(int64(500), nil).Once()

		id, err := srv.Save(m)

		assert.NoError(t, err)
		assert.Equal(t, int64(500), id)
		mock.AssertExpectationsForObjects(t, mRepo, sRepo)
	})

	t.Run("Sensor_Inactive_Error", func(t *testing.T) {
		srv, mRepo, sRepo := setupMetricServiceTest()
		m := models.Metric{SensorID: 1, Value: 25.5}

		// Датчик в статусі ERROR або INACTIVE
		sRepo.On("GetByID", 1).Return(models.Sensor{ID: 1, Status: models.StatusError}, nil).Once()

		id, err := srv.Save(m)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "is not active")
		assert.Equal(t, int64(0), id)
		// Репозиторій метрик НЕ повинен викликатися
		mRepo.AssertNotCalled(t, "Create", mock.Anything)
	})

	t.Run("Sensor_Not_Found", func(t *testing.T) {
		srv, _, sRepo := setupMetricServiceTest()
		m := models.Metric{SensorID: 999, Value: 10.0}

		sRepo.On("GetByID", 999).Return(models.Sensor{}, errors.New("not found")).Once()

		_, err := srv.Save(m)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "sensor 999 not found")
	})
}

func TestMetricService_GetHistory(t *testing.T) {
	t.Run("Invalid_Time_Range", func(t *testing.T) {
		srv, _, _ := setupMetricServiceTest()
		now := time.Now()

		// 'from' пізніше ніж 'to'
		_, err := srv.GetHistory(1, now, now.Add(-1*time.Hour))

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be after")
	})
}
