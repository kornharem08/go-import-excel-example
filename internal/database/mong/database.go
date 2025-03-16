package mong

import "go.mongodb.org/mongo-driver/mongo"

// IDatabase is interface of database structure.
type IDatabase interface {
	// Collection is get a collection by name.
	Collection(name string) ICollection
}

// Database is mongo.Database structure wrapper.
type Database struct {
	db *mongo.Database
}

func (database Database) Collection(name string) ICollection {
	return Collection{database.db.Collection(name)}
}
