package service

import (
	"testing"

	"agro-data-management-system/internal/models"
	"agro-data-management-system/internal/repository/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestCropService_Create(t *testing.T) {
	logger := zap.NewNop() // "Німий" логер для тестів

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mocks.CropRepository)
		srv := NewCropService(mockRepo, logger)

		crop := models.Crop{Name: "Пшениця", Variety: "Скарбниця"}
		mockRepo.On("Create", mock.AnythingOfType("models.Crop")).Return(1, nil).Once()

		id, err := srv.Create(crop)

		assert.NoError(t, err)
		assert.Equal(t, 1, id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Validation Error (Empty Name)", func(t *testing.T) {
		// Створюємо новий чистий мок ТУТ
		mockRepo := new(mocks.CropRepository)
		srv := NewCropService(mockRepo, logger)

		invalidCrop := models.Crop{Name: "", Variety: "Сорт"}

		id, err := srv.Create(invalidCrop)

		// Перевіряємо, що повернулася помилка валідації
		assert.Error(t, err)
		assert.Equal(t, 0, id)

		mockRepo.AssertNotCalled(t, "Create", mock.Anything)
	})
}
