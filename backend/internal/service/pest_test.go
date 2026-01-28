package service

import (
	"testing"

	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestPestService_Create(t *testing.T) {
	logger := zap.NewNop()

	t.Run("Success", func(t *testing.T) {
		// Створюємо мок ТУТ
		mockRepo := new(mocks.PestRepository)
		srv := NewPestService(mockRepo, logger)

		pest := models.Pest{Name: "Сарана", ScientificName: "Locusta migratoria"}
		expectedPest := models.Pest{ID: 1, Name: "Сарана", ScientificName: "Locusta migratoria"}
		mockRepo.On("Create", pest).Return(expectedPest, nil).Once()

		result, err := srv.Create(pest)

		assert.NoError(t, err)
		assert.Equal(t, expectedPest, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Validation_Error_Short_Latin_Name", func(t *testing.T) {
		// Створюємо НОВИЙ мок ТУТ
		mockRepo := new(mocks.PestRepository)
		srv := NewPestService(mockRepo, logger)

		// "Bug" — це 3 символи, а ми вимагаємо min=5
		invalidPest := models.Pest{Name: "Жук", ScientificName: "Bug"}

		result, err := srv.Create(invalidPest)

		// 1. Має бути помилка
		assert.Error(t, err)
		// 2. Повертається те саме невалідне значення
		assert.Equal(t, invalidPest, result)
		// 3. РЕПОЗИТОРІЙ НЕ МАЄ ВИКЛИКАТИСЯ
		mockRepo.AssertNotCalled(t, "Create", mock.Anything)
	})
}
