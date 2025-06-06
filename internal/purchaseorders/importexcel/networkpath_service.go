package importexcel

import (
	"purchase-record/internal/models"
)

type INetworkPathService interface {
	GetOrdersFromPath(filePath string) ([]models.PurchaseOrder, error)
}

type NetworkPathService struct {
	Repository INetworkPathRepository
}

func NewNetworkPathService() INetworkPathService {
	return &NetworkPathService{
		Repository: NewNetworkPathRepository(),
	}
}

func (s *NetworkPathService) GetOrdersFromPath(filePath string) ([]models.PurchaseOrder, error) {
	// Windows UNC paths with double backslashes are already correctly formatted
	// for the excelize library to process, so no conversion is needed

	// Get all orders from the repository
	return s.Repository.GetOrdersFromNetworkPath(filePath)
}
