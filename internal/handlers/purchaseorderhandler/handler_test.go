package purchaseorderhandler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"purchase-record/internal/models"
	"purchase-record/internal/purchaseorders/importexcel/mocks"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_GetOrdersFromNetworkPath(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockSetup      func(*mocks.INetworkPathService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:        "successful retrieval with no filter",
			queryParams: "",
			mockSetup: func(m *mocks.INetworkPathService) {
				m.On("GetOrdersFromPath", mock.Anything, "").
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
							PaymentTerm:        stringPtr("Net 30"),
							RequestDate:        stringPtr("2024-01-01"),
							DeliveryDate:       stringPtr("2024-06-30"),
							Status:             stringPtr("Active"),
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": []models.PurchaseOrder{
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
						PaymentTerm:        stringPtr("Net 30"),
						RequestDate:        stringPtr("2024-01-01"),
						DeliveryDate:       stringPtr("2024-06-30"),
						Status:             stringPtr("Active"),
					},
				},
			},
		},
		{
			name:        "successful retrieval with job_id_no filter",
			queryParams: "?job_id_no=123",
			mockSetup: func(m *mocks.INetworkPathService) {
				m.On("GetOrdersFromPath", mock.Anything, "123").
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
							PaymentTerm:        stringPtr("Net 30"),
							RequestDate:        stringPtr("2024-01-01"),
							DeliveryDate:       stringPtr("2024-06-30"),
							Status:             stringPtr("Active"),
						},
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"data": []models.PurchaseOrder{
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
						PaymentTerm:        stringPtr("Net 30"),
						RequestDate:        stringPtr("2024-01-01"),
						DeliveryDate:       stringPtr("2024-06-30"),
						Status:             stringPtr("Active"),
					},
				},
			},
		},
		{
			name:        "service error",
			queryParams: "",
			mockSetup: func(m *mocks.INetworkPathService) {
				m.On("GetOrdersFromPath", mock.Anything, "").
					Return(nil, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "service error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := mocks.NewINetworkPathService(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockService)
			}

			// Create handler with mock service
			handler := &Handler{
				NetworkPathService: mockService,
			}

			// Create test request
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/purchaseorders/import-network"+tt.queryParams, nil)

			// Call the handler
			handler.GetOrdersFromNetworkPath(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.JSONEq(t, toJSON(tt.expectedBody), w.Body.String())

			// Verify that all expected mock calls were made
			mockService.AssertExpectations(t)
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func toJSON(v interface{}) string {
	json, _ := json.Marshal(v)
	return string(json)
}
