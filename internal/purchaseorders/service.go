package purchaseorders

import (
	"context"
	"purchase-record/internal/database/mong"
	"purchase-record/internal/models"
)

type IService interface {
	CreatePurchaseOrder(ctx context.Context, order []models.PurchaseOrder) error
	GetPurchaseOrderList(ctx context.Context, query models.RequestQuery) ([]models.PurchaseOrder, int64, error)
}

type Service struct {
	Repository IRepository
}

func NewService(dbconn mong.IConnect) IService {
	return &Service{
		Repository: NewRepository(dbconn),
	}
}

func (s *Service) CreatePurchaseOrder(ctx context.Context, order []models.PurchaseOrder) error {
	return s.Repository.Create(ctx, order)
}

func (s *Service) GetPurchaseOrderList(ctx context.Context, query models.RequestQuery) ([]models.PurchaseOrder, int64, error) {
	return s.Repository.GetList(ctx, query)
}
