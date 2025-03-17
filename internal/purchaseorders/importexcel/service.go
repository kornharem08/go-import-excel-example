package importexcel

import (
	"mime/multipart"
	"purchase-record/internal/models"
)

type IService interface {
	GetOrders(file multipart.File) ([]models.PurchaseOrder, error)
}

type Service struct {
	Repository IRepository
}

func NewService() IService {
	return &Service{
		Repository: NewRepository(),
	}
}

func (s *Service) GetOrders(file multipart.File) ([]models.PurchaseOrder, error) {
	return s.Repository.GetOrdersFromExcel(file)
}
