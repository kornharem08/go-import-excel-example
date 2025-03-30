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
			tt.setupMock(mockService)

			// Create handler with mock service
			handler := &Handler{
				NetworkPathService: mockService,
			}

			// Create test context
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

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
