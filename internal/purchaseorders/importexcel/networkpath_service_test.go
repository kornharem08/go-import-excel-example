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
		jobIDNo        string
		purchasing     string
		status         string
		mockSetup      func(*mocks.INetworkPathRepository)
		expectedOrders []models.PurchaseOrder
		expectedError  error
	}{
		{
			name:       "successful retrieval with no filter",
			filePath:   "test.xlsx",
			jobIDNo:    "",
			purchasing: "",
			status:     "",
			mockSetup: func(m *mocks.INetworkPathRepository) {
				m.On("GetOrdersFromNetworkPath", "test.xlsx", "").
					Return([]models.PurchaseOrder{
						{
							JobIDNo:            stringPtr("123"),
							SalesTeam:          stringPtr("Team A"),
							ProjectManager:     stringPtr("John"),
							Purchasing:         stringPtr("Jane"),
							CustomerPO:         stringPtr("PO123"),
							JobAmount:          intPtr(1000),
							PeriodStart:        stringPtr("2024-01-01"),
							PeriodEnd:          stringPtr("2024-12-31"),
							Customer:           stringPtr("Customer A"),
							ProductCode:        stringPtr("PROD123"),
							ProductDescription: stringPtr("Test Product"),
							Ordered:            intPtr(100),
							Received:           intPtr(50),
							Remain:             intPtr(50),
							Currency:           stringPtr("USD"),
							UnitListPrice:      intPtr(10),
							ExtendListPrice:    intPtr(1000),
							DiscountPercent:    intPtr(5),
							DiscountAmount:     intPtr(50),
							ExtendUnitNetPrice: intPtr(9),
							ExtendNetPrice:     intPtr(900),
							DeliveryDate:       stringPtr("2024-06-30"),
							Status:             stringPtr("Active"),
						},
					}, nil)
			},
			expectedOrders: []models.PurchaseOrder{
				{
					JobIDNo:            stringPtr("123"),
					SalesTeam:          stringPtr("Team A"),
					ProjectManager:     stringPtr("John"),
					Purchasing:         stringPtr("Jane"),
					CustomerPO:         stringPtr("PO123"),
					JobAmount:          intPtr(1000),
					PeriodStart:        stringPtr("2024-01-01"),
					PeriodEnd:          stringPtr("2024-12-31"),
					Customer:           stringPtr("Customer A"),
					ProductCode:        stringPtr("PROD123"),
					ProductDescription: stringPtr("Test Product"),
					Ordered:            intPtr(100),
					Received:           intPtr(50),
					Remain:             intPtr(50),
					Currency:           stringPtr("USD"),
					UnitListPrice:      intPtr(10),
					ExtendListPrice:    intPtr(1000),
					DiscountPercent:    intPtr(5),
					DiscountAmount:     intPtr(50),
					ExtendUnitNetPrice: intPtr(9),
					ExtendNetPrice:     intPtr(900),
					DeliveryDate:       stringPtr("2024-06-30"),
					Status:             stringPtr("Active"),
				},
			},
			expectedError: nil,
		},
		{
			name:       "successful retrieval with job_id_no filter",
			filePath:   "test.xlsx",
			jobIDNo:    "123",
			purchasing: "",
			status:     "",
			mockSetup: func(m *mocks.INetworkPathRepository) {
				m.On("GetOrdersFromNetworkPath", "test.xlsx", "123").
					Return([]models.PurchaseOrder{
						{
							JobIDNo:            stringPtr("123"),
							SalesTeam:          stringPtr("Team A"),
							ProjectManager:     stringPtr("John"),
							Purchasing:         stringPtr("Jane"),
							CustomerPO:         stringPtr("PO123"),
							JobAmount:          intPtr(1000),
							PeriodStart:        stringPtr("2024-01-01"),
							PeriodEnd:          stringPtr("2024-12-31"),
							Customer:           stringPtr("Customer A"),
							ProductCode:        stringPtr("PROD123"),
							ProductDescription: stringPtr("Test Product"),
							Ordered:            intPtr(100),
							Received:           intPtr(50),
							Remain:             intPtr(50),
							Currency:           stringPtr("USD"),
							UnitListPrice:      intPtr(10),
							ExtendListPrice:    intPtr(1000),
							DiscountPercent:    intPtr(5),
							DiscountAmount:     intPtr(50),
							ExtendUnitNetPrice: intPtr(9),
							ExtendNetPrice:     intPtr(900),
							DeliveryDate:       stringPtr("2024-06-30"),
							Status:             stringPtr("Active"),
						},
					}, nil)
			},
			expectedOrders: []models.PurchaseOrder{
				{
					JobIDNo:            stringPtr("123"),
					SalesTeam:          stringPtr("Team A"),
					ProjectManager:     stringPtr("John"),
					Purchasing:         stringPtr("Jane"),
					CustomerPO:         stringPtr("PO123"),
					JobAmount:          intPtr(1000),
					PeriodStart:        stringPtr("2024-01-01"),
					PeriodEnd:          stringPtr("2024-12-31"),
					Customer:           stringPtr("Customer A"),
					ProductCode:        stringPtr("PROD123"),
					ProductDescription: stringPtr("Test Product"),
					Ordered:            intPtr(100),
					Received:           intPtr(50),
					Remain:             intPtr(50),
					Currency:           stringPtr("USD"),
					UnitListPrice:      intPtr(10),
					ExtendListPrice:    intPtr(1000),
					DiscountPercent:    intPtr(5),
					DiscountAmount:     intPtr(50),
					ExtendUnitNetPrice: intPtr(9),
					ExtendNetPrice:     intPtr(900),
					DeliveryDate:       stringPtr("2024-06-30"),
					Status:             stringPtr("Active"),
				},
			},
			expectedError: nil,
		},
		{
			name:       "successful retrieval with purchasing filter",
			filePath:   "test.xlsx",
			jobIDNo:    "",
			purchasing: "Jane",
			status:     "",
			mockSetup: func(m *mocks.INetworkPathRepository) {
				m.On("GetOrdersFromNetworkPath", "test.xlsx", "").
					Return([]models.PurchaseOrder{
						{
							JobIDNo:            stringPtr("123"),
							Purchasing:         stringPtr("Jane"),
							ProductDescription: stringPtr("Test Product"),
							Status:             stringPtr("Active"),
						},
					}, nil)
			},
			expectedOrders: []models.PurchaseOrder{
				{
					JobIDNo:            stringPtr("123"),
					Purchasing:         stringPtr("Jane"),
					ProductDescription: stringPtr("Test Product"),
					Status:             stringPtr("Active"),
				},
			},
			expectedError: nil,
		},
		{
			name:       "successful retrieval with status filter",
			filePath:   "test.xlsx",
			jobIDNo:    "",
			purchasing: "",
			status:     "Active",
			mockSetup: func(m *mocks.INetworkPathRepository) {
				m.On("GetOrdersFromNetworkPath", "test.xlsx", "").
					Return([]models.PurchaseOrder{
						{
							JobIDNo:            stringPtr("123"),
							Purchasing:         stringPtr("Jane"),
							ProductDescription: stringPtr("Test Product"),
							Status:             stringPtr("Active"),
						},
					}, nil)
			},
			expectedOrders: []models.PurchaseOrder{
				{
					JobIDNo:            stringPtr("123"),
					Purchasing:         stringPtr("Jane"),
					ProductDescription: stringPtr("Test Product"),
					Status:             stringPtr("Active"),
				},
			},
			expectedError: nil,
		},
		{
			name:       "successful retrieval with multiple filters",
			filePath:   "test.xlsx",
			jobIDNo:    "123",
			purchasing: "Jane",
			status:     "Active",
			mockSetup: func(m *mocks.INetworkPathRepository) {
				m.On("GetOrdersFromNetworkPath", "test.xlsx", "123").
					Return([]models.PurchaseOrder{
						{
							JobIDNo:            stringPtr("123"),
							Purchasing:         stringPtr("Jane"),
							ProductDescription: stringPtr("Test Product"),
							Status:             stringPtr("Active"),
						},
					}, nil)
			},
			expectedOrders: []models.PurchaseOrder{
				{
					JobIDNo:            stringPtr("123"),
					Purchasing:         stringPtr("Jane"),
					ProductDescription: stringPtr("Test Product"),
					Status:             stringPtr("Active"),
				},
			},
			expectedError: nil,
		},
		{
			name:       "repository error",
			filePath:   "test.xlsx",
			jobIDNo:    "",
			purchasing: "",
			status:     "",
			mockSetup: func(m *mocks.INetworkPathRepository) {
				m.On("GetOrdersFromNetworkPath", "test.xlsx", "").
					Return(nil, errors.New("repository error"))
			},
			expectedOrders: nil,
			expectedError:  errors.New("repository error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := mocks.NewINetworkPathRepository(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockRepo)
			}

			// Create service with mock repository
			service := &NetworkPathService{
				Repository: mockRepo,
			}

			// Call the method
			orders, err := service.GetOrdersFromPath(tt.filePath, tt.jobIDNo)

			// Assert results
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOrders, orders)

			// Verify that all expected mock calls were made
			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
