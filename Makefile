export MONGO_URI=mongodb://localhost:27017
export MONGODB_DATABASE_NAME=puchase_management
export MONGODB_PURCHASE_ORDER_COLLECTION_NAME=purchase_order

run-api:
	go run ./cmd/main.go

