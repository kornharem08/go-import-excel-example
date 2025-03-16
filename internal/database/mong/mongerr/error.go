package mongerr

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	// ErrConnectionNotResponse is connection ping but not pong responsed.
	ErrConnectionNotResponse = errors.New("connection is not pong")

	// ErrDatabaseNameIsEmpty is database name is empty.
	ErrDatabaseNameIsEmpty = errors.New("database name must be present")

	// ErrNoDocuments is returned by SingleResult methods when the operation that created the SingleResult did not return any documents.
	ErrNoDocuments = mongo.ErrNoDocuments
)

// IsErrNoDocuments check error is mongo.ErrNoDocuments.
func IsErrNoDocuments(err error) bool {
	return err == ErrNoDocuments
}

// IsErrDuplicateKey check duplicate key error.
func IsErrDuplicateKey(err error) bool {
	return mongo.IsDuplicateKeyError(err)
}
