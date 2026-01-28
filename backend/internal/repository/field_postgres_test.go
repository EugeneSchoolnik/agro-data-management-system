package repository

import (
	"agro-data-management-system/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldRepository_Lifecycle(t *testing.T) {
	// Очищаємо всі таблиці
	db := setupTestDB(t, "crops", "fields")
	defer db.Close()

	cropRepo := NewCropPostgres(db)
	fieldRepo := NewFieldPostgres(db)

	// 1. Створюємо залежність (культуру)
	crop, err := cropRepo.Create(models.Crop{
		Name:    "Пшениця",
		Variety: "Скарбниця",
	})
	assert.NoError(t, err)
	cropID := crop.ID

	// 2. Тестуємо Create Field
	newField := models.Field{
		Name:     "Північний участок",
		Area:     150.5,
		Location: "48.8584, 2.2945",
		CropID:   cropID,
	}

	field, err := fieldRepo.Create(newField)
	assert.NoError(t, err)
	assert.Greater(t, field.ID, 0)

	fieldID := field.ID

	// 3. Тестуємо GetByIDWithCrop (Перевірка JOIN)
	fieldWithCrop, err := fieldRepo.GetByIDWithCrop(fieldID)
	assert.NoError(t, err)
	assert.Equal(t, "Пшениця", fieldWithCrop.CropName)
	assert.Equal(t, "Північний участок", fieldWithCrop.Name)

	// 4. Тестуємо Update
	fieldWithCrop.Field.Name = "Оновлений участок"
	err = fieldRepo.Update(fieldWithCrop.Field)
	assert.NoError(t, err)

	updated, _ := fieldRepo.GetByIDWithCrop(fieldID)
	assert.Equal(t, "Оновлений участок", updated.Name)

	// 5. Тестуємо Delete
	err = fieldRepo.Delete(fieldID)
	assert.NoError(t, err)
}
