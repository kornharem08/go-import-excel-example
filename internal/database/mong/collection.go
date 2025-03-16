package mong

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ICollection is interface of collection structure.
type ICollection interface {
	// Find executes a find command and returns a Cursor over the matching documents in the collection.
	Find(ctx context.Context, result any, filter any, opts ...*options.FindOptions) error

	// FindOne executes a find command and returns a SingleResult for one document in the collection.
	FindOne(ctx context.Context, result any, filter any, opts ...*options.FindOneOptions) error

	// InsertOne executes an insert command to insert a single document into the collection.
	InsertOne(ctx context.Context, document any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)

	// InsertMany executes an insert command to insert multiple documents into the collection.
	// If write errors occur during the operation (e.g. duplicate key error), this method returns a BulkWriteException error.
	InsertMany(ctx context.Context, documents []any, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error)

	// UpdateByID executes an update command to update the document whose _id value matches the provided ID in the collection.
	// This is equivalent to running UpdateOne(ctx, bson.D{{"_id", id}}, update, opts...).
	UpdateByID(ctx context.Context, id any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)

	// UpdateOne executes an update command to update at most one document in the collection.
	UpdateOne(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)

	// UpdateMany executes an update command to update documents in the collection.
	UpdateMany(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error)

	// ReplaceOne executes an update command to replace at most one document in the collection.
	//
	// The filter parameter must be a document containing query operators and can be used to select the document to be
	// replaced. It cannot be nil. If the filter does not match any documents, the operation will succeed and an
	// UpdateResult with a MatchedCount of 0 will be returned. If the filter matches multiple documents, one will be
	// selected from the matched set and MatchedCount will equal 1.
	//
	// The replacement parameter must be a document that will be used to replace the selected document. It cannot be nil
	// and cannot contain any update operators (https://www.mongodb.com/docs/manual/reference/operator/update/).
	//
	// The opts parameter can be used to specify options for the operation (see the options.ReplaceOptions documentation).
	//
	// For more information about the command, see https://www.mongodb.com/docs/manual/reference/command/update/.
	ReplaceOne(ctx context.Context, filter any, replacement any, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error)

	// FindOneAndReplace executes a findAndModify command to replace at most one document in the collection and returns the document as it appeared after replacement.
	FindOneAndReplace(ctx context.Context, result any, filter any, replacement any, opts ...*options.FindOneAndReplaceOptions) error

	// DeleteOne executes a delete command to delete at most one document from the collection.
	DeleteOne(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)

	// DeleteMany executes a delete command to delete documents from the collection.
	DeleteMany(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error)

	// Aggregate executes an aggregate command against the collection and returns a cursor over the resulting documents.
	Aggregate(ctx context.Context, result any, pipeline any, opts ...*options.AggregateOptions) error

	// CountDocuments returns the number of documents in the collection.
	// For a fast count of the documents in the collection, see the EstimatedDocumentCount method.
	CountDocuments(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error)

	// FindOneAndUpdate executes a findAndModify command to update at most one document in the collection and returns the document as it appeared after updating.
	FindOneAndUpdate(ctx context.Context, result any, filter any, update any, opts ...*options.FindOneAndUpdateOptions) error

	// BulkWrite performs a bulk write operation (https://www.mongodb.com/docs/manual/core/bulk-write-operations/).
	BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error)

	// Drop drops the collection on the server. This method ignores "namespace not found" errors so it is safe to drop a collection that does not exist on the server.
	Drop(ctx context.Context) error
}

// Collection is mongo.Collection structure wrapper.
type Collection struct {
	col *mongo.Collection
}

func (collection Collection) Find(ctx context.Context, result any, filter any, opts ...*options.FindOptions) error {
	cursor, err := collection.col.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}

	if err := cursor.All(ctx, result); err != nil {
		return err
	}

	return nil
}

func (collection Collection) FindOne(ctx context.Context, result any, filter any, opts ...*options.FindOneOptions) error {
	response := collection.col.FindOne(ctx, filter, opts...)

	if err := response.Err(); err != nil {
		return err
	}

	if err := response.Decode(result); err != nil {
		return err
	}

	return nil
}

func (collection Collection) InsertOne(ctx context.Context, document any, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	result, err := collection.col.InsertOne(ctx, document, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (collection Collection) InsertMany(ctx context.Context, documents []any, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	result, err := collection.col.InsertMany(ctx, documents, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (collection Collection) UpdateByID(ctx context.Context, id any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	result, err := collection.col.UpdateByID(ctx, id, update, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (collection Collection) UpdateOne(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	result, err := collection.col.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (collection Collection) UpdateMany(ctx context.Context, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	result, err := collection.col.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (collection Collection) ReplaceOne(ctx context.Context, filter any, replacement any, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	return collection.col.ReplaceOne(ctx, filter, replacement, opts...)
}

func (collection Collection) FindOneAndReplace(ctx context.Context, result any, filter any, replacement any, opts ...*options.FindOneAndReplaceOptions) error {
	opt := options.FindOneAndReplace()
	opt.SetReturnDocument(options.After)

	opts = append(opts, opt)

	if err := collection.col.FindOneAndReplace(ctx, filter, replacement, opts...).Decode(result); err != nil {
		return err
	}

	return nil
}

func (collection Collection) DeleteOne(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	result, err := collection.col.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (collection Collection) DeleteMany(ctx context.Context, filter any, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	result, err := collection.col.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (collection Collection) Aggregate(ctx context.Context, result any, pipeline interface{}, opts ...*options.AggregateOptions) error {
	cursor, err := collection.col.Aggregate(ctx, pipeline, opts...)
	if err != nil {
		return err
	}

	if err := cursor.All(ctx, result); err != nil {
		return err
	}

	return nil
}

func (collection Collection) CountDocuments(ctx context.Context, filter any, opts ...*options.CountOptions) (int64, error) {
	return collection.col.CountDocuments(ctx, filter, opts...)
}

func (collection Collection) FindOneAndUpdate(ctx context.Context, result any, filter any, update any, opts ...*options.FindOneAndUpdateOptions) error {
	opt := options.FindOneAndUpdate()
	opt.SetReturnDocument(options.After)

	opts = append(opts, opt)

	if err := collection.col.FindOneAndUpdate(ctx, filter, update, opts...).Decode(result); err != nil {
		return err
	}

	return nil
}

func (collection Collection) BulkWrite(ctx context.Context, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	return collection.col.BulkWrite(ctx, models, opts...)
}

func (collection Collection) Drop(ctx context.Context) error {
	return collection.col.Drop(ctx)
}
