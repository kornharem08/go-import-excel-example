package importexcel_test

import (
	"errors"
	"mime/multipart"
	"purchase-record/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"

	// Assuming mockery generated mock in mocks folder
	"purchase-record/internal/purchaseorders/importexcel"
	"purchase-record/internal/purchaseorders/importexcel/mocks"
	"purchase-record/internal/purchaseorders/utils"

	"github.com/stretchr/testify/mock"
)

// Test case 1: Successful order retrieval
func TestService_GetOrders_Success(t *testing.T) {
	var (
		mockRepo *mocks.IRepository
		service  importexcel.IService
	)
	mockRepo = new(mocks.IRepository)

	expectedOrders := []models.PurchaseOrder{
		{
			JobIDNo:            utils.StringOrNil("1"),
			ProductDescription: utils.StringOrNil("Test Product"),
			Status:             utils.StringOrNil("Pending"),
		},
	}

	mockRepo.On("GetOrdersFromExcel", mock.Anything).Return(expectedOrders, nil)

	var mockFile multipart.File
	service = &importexcel.Service{
		Repository: mockRepo,
	}
	orders, err := service.GetOrders(mockFile)

	assert.NoError(t, err)
	assert.Equal(t, expectedOrders, orders)
	mockRepo.AssertExpectations(t)
}

// Test case 2: Repository returns error
func TestService_GetOrders_RepositoryError(t *testing.T) {
	var (
		mockRepo *mocks.IRepository
		service  importexcel.IService
	)
	mockRepo = new(mocks.IRepository)
	expectedError := errors.New("repository error")
	mockRepo.On("GetOrdersFromExcel", mock.Anything).Return(nil, expectedError)

	var mockFile multipart.File
	service = &importexcel.Service{
		Repository: mockRepo,
	}
	result, err := service.GetOrders(mockFile)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, expectedError, err)
	assert.Nil(t, result)
}

// // Test case 3: Empty file returns empty orders
func TestService_GetOrders_EmptyFile(t *testing.T) {
	var (
		mockRepo *mocks.IRepository
		service  importexcel.IService
	)
	mockRepo = new(mocks.IRepository)

	expectedOrders := []models.PurchaseOrder{}

	mockRepo.On("GetOrdersFromExcel", mock.Anything).Return(expectedOrders, nil)

	var mockFile multipart.File
	service = &importexcel.Service{
		Repository: mockRepo,
	}
	result, err := service.GetOrders(mockFile)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, expectedOrders, result)
}
