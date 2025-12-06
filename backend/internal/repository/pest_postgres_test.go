package repository

import (
	"testing"
	"agro-data-management-system/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestPestRepository_Lifecycle(t *testing.T) {
	db := setupTestDB(t, "pests")
	defer db.Close()

	repo := NewPestPostgres(db)

	// 1. Create
	newPest := models.Pest{
		Name:           "Клоп шкідлива черепашка",
		ScientificName: "Eurygaster integriceps",
	}
	id, err := repo.Create(newPest)
	assert.NoError(t, err)
	assert.Greater(t, id, 0)

	// 2. GetByID
	found, err := repo.GetByID(id)
	assert.NoError(t, err)
	assert.Equal(t, newPest.ScientificName, found.ScientificName)

	// 3. Update
	found.Name = "Оновлена назва"
	err = repo.Update(found)
	assert.NoError(t, err)

	updated, _ := repo.GetByID(id)
	assert.Equal(t, "Оновлена назва", updated.Name)

	// 4. GetAll
	all, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, all, 1)

	// 5. Delete
	err = repo.Delete(id)
	assert.NoError(t, err)
}