package importexcel

import (
	"errors"
	"purchase-record/internal/models"
	"purchase-record/internal/purchaseorders/importexcel/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNetworkPathService_GetOrdersFromPath(t *testing.T) {
	tests := []struct {
		name           string
		filePath       string
		mockSetup      func(*mocks.INetworkPathRepository)
		expectedOrders []models.PurchaseOrder
		expectedError  error
	}{
		{
			name:     "successful retrieval of all orders",
			filePath: "test.xlsx",
			mockSetup: func(m *mocks.INetworkPathRepository) {
				m.On("GetOrdersFromNetworkPath", "test.xlsx").
					Return([]models.PurchaseOrder{
						{
							JobIDNo:            stringPtr("123"),
							Type:               stringPtr("Standard"),
							SalesTeam:          stringPtr("Team A"),
							ProjectManager:     stringPtr("John"),
							Customer:           stringPtr("Customer A"),
							ProductCode:        stringPtr("PROD123"),
							ProductDescription: stringPtr("Test Product"),
							Ordered:            intPtr(100),
							Received:           intPtr(50),
							Remain:             intPtr(50),
							PR:                 stringPtr("PR123"),
							PRDate:             stringPtr("2024-01-01"),
							PO:                 stringPtr("PO123"),
							PODate:             stringPtr("2024-01-15"),
							Distribution:       stringPtr("Dept1"),
							RequestDate:        stringPtr("2024-01-01"),
							DeliveryDate:       stringPtr("2024-06-30"),
							Status:             stringPtr("Active"),
						},
						{
							JobIDNo:            stringPtr("456"),
							Type:               stringPtr("Express"),
							SalesTeam:          stringPtr("Team B"),
							ProjectManager:     stringPtr("Jane"),
							Customer:           stringPtr("Customer B"),
							ProductCode:        stringPtr("PROD456"),
							ProductDescription: stringPtr("Test Product 2"),
							Ordered:            intPtr(200),
							Received:           intPtr(100),
							Remain:             intPtr(100),
							PR:                 stringPtr("PR456"),
							PRDate:             stringPtr("2024-02-01"),
							PO:                 stringPtr("PO456"),
							PODate:             stringPtr("2024-02-15"),
							Distribution:       stringPtr("Dept2"),
							RequestDate:        stringPtr("2024-02-01"),
							DeliveryDate:       stringPtr("2024-07-30"),
							Status:             stringPtr("Active"),
						},
					}, nil)
			},
			expectedOrders: []models.PurchaseOrder{
				{
					JobIDNo:            stringPtr("123"),
					Type:               stringPtr("Standard"),
					SalesTeam:          stringPtr("Team A"),
					ProjectManager:     stringPtr("John"),
					Customer:           stringPtr("Customer A"),
					ProductCode:        stringPtr("PROD123"),
					ProductDescription: stringPtr("Test Product"),
					Ordered:            intPtr(100),
					Received:           intPtr(50),
					Remain:             intPtr(50),
					PR:                 stringPtr("PR123"),
					PRDate:             stringPtr("2024-01-01"),
					PO:                 stringPtr("PO123"),
					PODate:             stringPtr("2024-01-15"),
					Distribution:       stringPtr("Dept1"),
					RequestDate:        stringPtr("2024-01-01"),
					DeliveryDate:       stringPtr("2024-06-30"),
					Status:             stringPtr("Active"),
				},
				{
					JobIDNo:            stringPtr("456"),
					Type:               stringPtr("Express"),
					SalesTeam:          stringPtr("Team B"),
					ProjectManager:     stringPtr("Jane"),
					Customer:           stringPtr("Customer B"),
					ProductCode:        stringPtr("PROD456"),
					ProductDescription: stringPtr("Test Product 2"),
					Ordered:            intPtr(200),
					Received:           intPtr(100),
					Remain:             intPtr(100),
					PR:                 stringPtr("PR456"),
					PRDate:             stringPtr("2024-02-01"),
					PO:                 stringPtr("PO456"),
					PODate:             stringPtr("2024-02-15"),
					Distribution:       stringPtr("Dept2"),
					RequestDate:        stringPtr("2024-02-01"),
					DeliveryDate:       stringPtr("2024-07-30"),
					Status:             stringPtr("Active"),
				},
			},
			expectedError: nil,
		},
		{
			name:     "successful retrieval with empty result",
			filePath: "empty.xlsx",
			mockSetup: func(m *mocks.INetworkPathRepository) {
				m.On("GetOrdersFromNetworkPath", "empty.xlsx").
					Return([]models.PurchaseOrder{}, nil)
			},
			expectedOrders: []models.PurchaseOrder{},
			expectedError:  nil,
		},
		{
			name:     "repository error",
			filePath: "test.xlsx",
			mockSetup: func(m *mocks.INetworkPathRepository) {
				m.On("GetOrdersFromNetworkPath", "test.xlsx").
					Return(nil, errors.New("repository error"))
			},
			expectedOrders: nil,
			expectedError:  errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup the mock repository
			mockRepo := new(mocks.INetworkPathRepository)
			tt.mockSetup(mockRepo)

			// Create service with the mock repository
			service := &NetworkPathService{Repository: mockRepo}

			// Execute the test
			orders, err := service.GetOrdersFromPath(tt.filePath)

			// Verify results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, orders)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedOrders, orders)
			}

			// Verify that all expected calls were made
			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper functions to create pointers for string and int values
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
