package importexcel

import (
	"fmt"
	"os"
	"path/filepath"
	"purchase-record/internal/models"
	"purchase-record/internal/purchaseorders/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkPathRepository_GetOrdersFromNetworkPath(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "test-excel-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up temp directory after tests

	// Define test orders
	testOrder1 := models.PurchaseOrder{
		JobIDNo:            utils.StringOrNil("123"),
		Type:               utils.StringOrNil("Standard"),
		SalesTeam:          utils.StringOrNil("Team A"),
		ProjectManager:     utils.StringOrNil("John"),
		Customer:           utils.StringOrNil("Customer A"),
		ProductCode:        utils.StringOrNil("PROD123"),
		ProductDescription: utils.StringOrNil("Test Product"),
		Ordered:            utils.IntOrNil("100"),
		Received:           utils.IntOrNil("50"),
		Remain:             utils.IntOrNil("50"),
		PR:                 utils.StringOrNil("PR123"),
		PRDate:             utils.StringOrNil("2024-01-01"),
		PO:                 utils.StringOrNil("PO123"),
		PODate:             utils.StringOrNil("2024-01-15"),
		Distribution:       utils.StringOrNil("Dept1"),
		RequestDate:        utils.StringOrNil("2024-01-01"),
		DeliveryDate:       utils.StringOrNil("2024-06-30"),
		Status:             utils.StringOrNil("Active"),
	}

	testOrder2 := models.PurchaseOrder{
		JobIDNo:            utils.StringOrNil("456"),
		Type:               utils.StringOrNil("Express"),
		SalesTeam:          utils.StringOrNil("Team B"),
		ProjectManager:     utils.StringOrNil("Jane"),
		Customer:           utils.StringOrNil("Customer B"),
		ProductCode:        utils.StringOrNil("PROD456"),
		ProductDescription: utils.StringOrNil("Test Product 2"),
		Ordered:            utils.IntOrNil("200"),
		Received:           utils.IntOrNil("100"),
		Remain:             utils.IntOrNil("100"),
		PR:                 utils.StringOrNil("PR456"),
		PRDate:             utils.StringOrNil("2024-02-01"),
		PO:                 utils.StringOrNil("PO456"),
		PODate:             utils.StringOrNil("2024-02-15"),
		Distribution:       utils.StringOrNil("Dept2"),
		RequestDate:        utils.StringOrNil("2024-02-01"),
		DeliveryDate:       utils.StringOrNil("2024-07-30"),
		Status:             utils.StringOrNil("Active"),
	}

	// Create test cases for the mock repository
	testCases := map[string]struct {
		orders []models.PurchaseOrder
		err    error
	}{
		"test1.xlsx": {
			orders: []models.PurchaseOrder{testOrder1},
			err:    nil,
		},
		"test2.xlsx": {
			orders: []models.PurchaseOrder{testOrder1, testOrder2},
			err:    nil,
		},
		"nonexistent.xlsx": {
			orders: nil,
			err:    fmt.Errorf("failed to open Excel file"),
		},
		"invalid.xlsx": {
			orders: nil,
			err:    fmt.Errorf("failed to open Excel file"),
		},
		"empty.xlsx": {
			orders: []models.PurchaseOrder{},
			err:    nil,
		},
	}

	// Create mock repository
	mockRepo := NewMockNetworkPathRepository(testCases)

	tests := []struct {
		name           string
		filePath       string
		expectedOrders []models.PurchaseOrder
		expectedError  error
	}{
		{
			name:           "successful retrieval from file with single order",
			filePath:       filepath.Join(tempDir, "test1.xlsx"),
			expectedOrders: []models.PurchaseOrder{testOrder1},
			expectedError:  nil,
		},
		{
			name:           "successful retrieval from file with multiple orders",
			filePath:       filepath.Join(tempDir, "test2.xlsx"),
			expectedOrders: []models.PurchaseOrder{testOrder1, testOrder2},
			expectedError:  nil,
		},
		{
			name:           "file not found error",
			filePath:       filepath.Join(tempDir, "nonexistent.xlsx"),
			expectedOrders: nil,
			expectedError:  fmt.Errorf("failed to open Excel file"),
		},
		{
			name:           "invalid excel file format",
			filePath:       filepath.Join(tempDir, "invalid.xlsx"),
			expectedOrders: nil,
			expectedError:  fmt.Errorf("failed to open Excel file"),
		},
		{
			name:           "empty excel file",
			filePath:       filepath.Join(tempDir, "empty.xlsx"),
			expectedOrders: []models.PurchaseOrder{},
			expectedError:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the method using mock repository
			orders, err := mockRepo.GetOrdersFromNetworkPath(tt.filePath)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOrders, orders)
		})
	}
}

// MockNetworkPathRepository is a mock implementation of INetworkPathRepository
type MockNetworkPathRepository struct {
	testCases map[string]struct {
		orders []models.PurchaseOrder
		err    error
	}
}

func NewMockNetworkPathRepository(testCases map[string]struct {
	orders []models.PurchaseOrder
	err    error
}) INetworkPathRepository {
	return &MockNetworkPathRepository{testCases: testCases}
}

func (m *MockNetworkPathRepository) GetOrdersFromNetworkPath(filePath string) ([]models.PurchaseOrder, error) {
	// Extract test case key from filePath - last part of path
	key := filepath.Base(filePath)
	if tc, ok := m.testCases[key]; ok {
		return tc.orders, tc.err
	}
	return nil, fmt.Errorf("no test case for file: %s", filePath)
}
