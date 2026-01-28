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

// helper функція для ініціалізації сервісу з моками
func setupFieldServiceTest() (*fieldService, *mocks.FieldRepository, *mocks.CropRepository) {
	fRepo := new(mocks.FieldRepository)
	cRepo := new(mocks.CropRepository)
	log := zap.NewNop()
	srv := NewFieldService(fRepo, cRepo, log).(*fieldService)
	return srv, fRepo, cRepo
}

func TestFieldService_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv, fRepo, cRepo := setupFieldServiceTest()
		field := models.Field{Name: "Південний сектор", Area: 100.5, Location: "45.0, 31.0", CropID: 1}

		cRepo.On("GetByID", 1).Return(models.Crop{ID: 1}, nil).Once()
		expectedField := models.Field{ID: 5, Name: "Південний сектор", Area: 100.5, Location: "45.0, 31.0", CropID: 1}
		fRepo.On("Create", field).Return(expectedField, nil).Once()

		result, err := srv.Create(field)

		assert.NoError(t, err)
		assert.Equal(t, expectedField, result)
		mock.AssertExpectationsForObjects(t, fRepo, cRepo)
	})

	t.Run("Validation_Error_Empty_Name", func(t *testing.T) {
		srv, fRepo, _ := setupFieldServiceTest()
		invalidField := models.Field{Name: "", Area: 10, Location: "0,0", CropID: 1}

		result, err := srv.Create(invalidField)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Field.Name")
		assert.Equal(t, invalidField, result)
		fRepo.AssertNotCalled(t, "Create", mock.Anything)
	})

	t.Run("Validation_Error_Negative_Area", func(t *testing.T) {
		srv, _, _ := setupFieldServiceTest()
		invalidField := models.Field{Name: "Тест", Area: -5.0, Location: "0,0", CropID: 1}

		_, err := srv.Create(invalidField)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Field.Area")
	})

	t.Run("Crop_Not_Found", func(t *testing.T) {
		srv, _, cRepo := setupFieldServiceTest()
		field := models.Field{Name: "Поле 1", Area: 10, Location: "0,0", CropID: 99}

		cRepo.On("GetByID", 99).Return(models.Crop{}, errors.New("not found")).Once()

		result, err := srv.Create(field)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "crop with id 99 not found")
		assert.Equal(t, field, result)
	})

	t.Run("Database_Error_On_Create", func(t *testing.T) {
		srv, fRepo, cRepo := setupFieldServiceTest()
		field := models.Field{Name: "Поле 1", Area: 10, Location: "0,0", CropID: 1}

		cRepo.On("GetByID", 1).Return(models.Crop{ID: 1}, nil).Once()
		fRepo.On("Create", field).Return(models.Field{}, errors.New("db connection lost")).Once()

		_, err := srv.Create(field)

		assert.Error(t, err)
		assert.Equal(t, "db connection lost", err.Error())
	})
}

func TestFieldService_GetByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv, fRepo, _ := setupFieldServiceTest()
		expected := models.FieldWithCrop{
			Field:    models.Field{ID: 1, Name: "Поле 1"},
			CropName: "Пшениця",
		}

		fRepo.On("GetByIDWithCrop", 1).Return(expected, nil).Once()

		res, err := srv.GetByID(1)

		assert.NoError(t, err)
		assert.Equal(t, "Пшениця", res.CropName)
		fRepo.AssertExpectations(t)
	})

	t.Run("Not_Found", func(t *testing.T) {
		srv, fRepo, _ := setupFieldServiceTest()
		fRepo.On("GetByIDWithCrop", 404).Return(models.FieldWithCrop{}, errors.New("not found")).Once()

		_, err := srv.GetByID(404)

		assert.Error(t, err)
	})
}

func TestFieldService_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		srv, fRepo, cRepo := setupFieldServiceTest()
		field := models.Field{ID: 1, Name: "Оновлена назва", Area: 20, Location: "1,1", CropID: 2}

		cRepo.On("GetByID", 2).Return(models.Crop{ID: 2}, nil).Once()
		fRepo.On("Update", field).Return(nil).Once()

		err := srv.Update(field)

		assert.NoError(t, err)
	})

	t.Run("Update_With_Invalid_Crop", func(t *testing.T) {
		srv, fRepo, cRepo := setupFieldServiceTest()
		field := models.Field{ID: 1, Name: "Тест", Area: 20, Location: "1,1", CropID: 999}

		cRepo.On("GetByID", 999).Return(models.Crop{}, errors.New("not found")).Once()

		err := srv.Update(field)

		assert.Error(t, err)
		fRepo.AssertNotCalled(t, "Update", mock.Anything)
	})
}
