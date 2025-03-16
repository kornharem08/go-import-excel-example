package config

type Config struct {
	// MongoDB
	URI             string `env:"MONGO_URI" default:"mongodb://localhost:27017"`
	MongoDBDatabase string `env:"MONGODB_DATABASE_NAME" default:"puchase_management"`
	PurchaseOrder   string `env:"MONGODB_PURCHASE_ORDER_COLLECTION_NAME" default:"purchase_order"`
}
