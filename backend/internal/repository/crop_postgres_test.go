package repository

import (
	"agro-data-management-system/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCropRepository_Lifecycle(t *testing.T) {
	// Викликаємо спільний хелпер, вказуючи, яку таблицю чистити
	db := setupTestDB(t, "crops")
	defer db.Close()

	repo := NewCropPostgres(db)

	newCrop := models.Crop{
		Name:        "Пшениця озима",
		Variety:     "Скарбниця",
		Description: "Тестовий опис для диплома",
	}

	// Create
	id, err := repo.Create(newCrop)
	assert.NoError(t, err)
	assert.Greater(t, id, 0)

	// GetByID
	found, err := repo.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, newCrop.Name, found.Name)

	// ... і так далі
}
