package purchaseorders

import (
	"context"
	"purchase-record/internal/config"
	"purchase-record/internal/database/environ"
	"purchase-record/internal/database/mong"
	"purchase-record/internal/models"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IRepository defines the interface for purchase order repository
type IRepository interface {
	Create(ctx context.Context, orders []models.PurchaseOrder) error
	GetList(ctx context.Context, query models.RequestQuery) ([]models.PurchaseOrder, int64, error)
}

// MongoRepository implements IRepository using MongoDB
type Repository struct {
	Collection mong.ICollection
}

// NewRepository creates a new MongoRepository instance
func NewRepository(dbconn mong.IConnect) IRepository {
	return &Repository{
		Collection: dbconn.Database().Collection(environ.Load[config.Config]().PurchaseOrder),
	}
}

// CreatePurchaseOrder inserts multiple purchase orders into MongoDB
func (r *Repository) Create(ctx context.Context, orders []models.PurchaseOrder) error {
	var data []any
	for _, order := range orders {
		data = append(data, order)
	}
	_, err := r.Collection.InsertMany(ctx, data)
	return err
}

func (r *Repository) GetList(ctx context.Context, query models.RequestQuery) ([]models.PurchaseOrder, int64, error) {
	var result []models.PurchaseOrder
	filter := bson.M{}
	// Create options for pagination
	findOptions := options.Find()
	findOptions.SetSkip((query.PageNo - 1) * query.PageSize)
	findOptions.SetLimit(query.PageSize)

	total, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to count purchase orders")
	}

	if err := r.Collection.Find(ctx, &result, filter, findOptions); err != nil {
		return []models.PurchaseOrder{}, 0, errors.Wrap(err, "failed to get purchase orders")
	}
	if len(result) == 0 {
		return []models.PurchaseOrder{}, 0, nil
	}
	return result, total, nil
}
