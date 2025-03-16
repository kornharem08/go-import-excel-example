package mong

import (
	"context"
	"log"
	"net/url"
	"purchase-record/internal/database/environ"
	"purchase-record/internal/database/mong/internal"
	"purchase-record/internal/database/mong/mongerr"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

// IConnect is interface of connection database.
type IConnect interface {
	// Close is disconnect a connection.
	Close() error

	// Database is interface of database wrapper.
	Database() IDatabase

	// NewSession is create a database transaction.
	NewSession() (ISession, error)
}

// Connect is structure for connect database.
type Connect struct {
	client   *mongo.Client
	database IDatabase
}

// getURI is get uri from environment variables under name of MONGO_URI,
// Validate it and add auth source if connect mode is not direct.
func getURI(databaseName string) (string, error) {
	// Get configurations
	conf := environ.Load[internal.Config]()

	// Check query parameters
	index := strings.LastIndex(conf.URI, "?")
	if index < 0 {
		return conf.URI, nil
	}

	// Parse uri
	parsed, err := url.Parse(conf.URI[index:])
	if err != nil {
		return "", err
	}

	// Get query parameters
	query := parsed.Query()

	// Add auth source if connect mode is direct
	if query.Get("connect") != connstring.SingleConnect.String() {
		query.Set("authSource", databaseName)
	}

	// Parse new query parameter
	uri := conf.URI[:index] + "?" + query.Encode()
	log.Fatal(uri)
	return uri, nil
}

// New database connection. Will connecting by URI and database name.
// This URI get from environment variables by MONGO_URI variable.
// The MONGO_URI example format is mongodb://<username>:<password>@<cluster_separate_with_comma>/?authSource=<database_name>&replicaSet=<replica_set>
func New(databaseName string, opts ...*options.DatabaseOptions) (IConnect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), (10 * time.Second))
	defer cancel()

	// Check database name
	if databaseName == "" {
		return nil, mongerr.ErrDatabaseNameIsEmpty
	}

	// Get URI
	uri, err := getURI(databaseName)
	if err != nil {
		return nil, err
	}

	// connect to database server
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	// Ping connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, mongerr.ErrConnectionNotResponse
	}

	return &Connect{
		// Store client
		client: client,

		// Choose database
		database: &Database{
			client.Database(databaseName, opts...),
		},
	}, nil
}

func (conn Connect) Close() error {
	if conn.client == nil {
		return nil
	}

	return conn.client.Disconnect(context.Background())
}

func (conn Connect) Database() IDatabase {
	return conn.database
}

func (conn Connect) NewSession() (ISession, error) {
	option := &options.SessionOptions{
		DefaultReadConcern:    readconcern.Snapshot(),
		DefaultWriteConcern:   writeconcern.Majority(),
		DefaultReadPreference: readpref.Primary(),
	}

	session, err := conn.client.StartSession(option)
	if err != nil {
		return nil, err
	}

	return Session{
		ctx:     context.Background(),
		session: session,
	}, nil
}
