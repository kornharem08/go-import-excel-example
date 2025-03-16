package importexcel

import (
	"mime/multipart"
	"purchase-record/internal/models"
)

type IService interface {
	GetOrders(file multipart.File) ([]models.PurchaseOrder, error)
}

type Service struct {
	repo IRepository
}

func NewService() IService {
	return &Service{
		repo: NewRepository(),
	}
}

func (s *Service) GetOrders(file multipart.File) ([]models.PurchaseOrder, error) {
	return s.repo.GetOrdersFromExcel(file)
}
