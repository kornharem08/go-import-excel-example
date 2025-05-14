package purchaseorderhandler

import (
	"net/http/httptest"
	"os"
	"path/filepath"
	"purchase-record/internal/models"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNetworkPathService is a mock implementation of INetworkPathService
type MockNetworkPathService struct {
	mock.Mock
}

// MockSettingPathService is a mock implementation of ISettingPathService
type MockSettingPathService struct {
	mock.Mock
}

func (m *MockSettingPathService) GetSettingPath(filePath string) ([]models.SettingExcelData, error) {
	args := m.Called(filePath)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.SettingExcelData), args.Error(1)
}

func (m *MockNetworkPathService) GetOrdersFromPath(filePath string, jobIDNo string) ([]models.PurchaseOrder, error) {
	args := m.Called(filePath, jobIDNo)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PurchaseOrder), args.Error(1)
}

func TestGetOrdersFromNetworkPath(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFilePath := filepath.Join(tempDir, "test.xlsx")

	// Create a test file
	err := os.WriteFile(testFilePath, []byte("test data"), 0644)
	assert.NoError(t, err)

	tests := []struct {
		name           string
		setupMock      func(*MockNetworkPathService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful read from original file",
			setupMock: func(m *MockNetworkPathService) {
				m.On("GetOrdersFromPath", mock.Anything, "").Return([]models.PurchaseOrder{}, nil)
			},
			expectedStatus: 200,
		},
		{
			name: "service error",
			setupMock: func(m *MockNetworkPathService) {
				m.On("GetOrdersFromPath", mock.Anything, "").Return(nil, assert.AnError)
			},
			expectedStatus: 500,
			expectedError:  assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockNetworkPathService)
			mockSettingService := new(MockSettingPathService)
			tt.setupMock(mockService)

			// Create handler with mock service
			handler := &Handler{
				NetworkPathService: mockService,
				SettingPathService: mockSettingService,
			}

			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set a test path in the request params
			c.Request = httptest.NewRequest("POST", "/purchaseorders", nil)
			c.Params = gin.Params{
				{Key: "path", Value: testFilePath},
			}

			// Add path to query string
			query := c.Request.URL.Query()
			query.Add("path", testFilePath)
			c.Request.URL.RawQuery = query.Encode()

			// Call the handler
			handler.GetOrdersFromNetworkPath(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}

func TestNewHandler(t *testing.T) {
	handler := NewHandler()
	assert.NotNil(t, handler)
}

func TestGetSettingPath(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockSettingPathService)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful retrieval",
			setupMock: func(m *MockSettingPathService) {
				m.On("GetSettingPath", mock.Anything).Return([]models.SettingExcelData{
					{Path: "test/path", Name: "Test"},
				}, nil)
			},
			expectedStatus: 200,
		},
		{
			name: "service error",
			setupMock: func(m *MockSettingPathService) {
				m.On("GetSettingPath", mock.Anything).Return(nil, assert.AnError)
			},
			expectedStatus: 500,
			expectedError:  "Failed to get setting path: " + assert.AnError.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock service
			mockService := new(MockSettingPathService)
			tt.setupMock(mockService)

			// Create handler with mock service
			handler := &Handler{
				SettingPathService: mockService,
			}

			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest("GET", "/purchaseorders/setting", nil)

			// Call the handler
			handler.GetSettingPath(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedError != "" {
				assert.Contains(t, w.Body.String(), tt.expectedError)
			}

			// Verify mock expectations
			mockService.AssertExpectations(t)
		})
	}
}
