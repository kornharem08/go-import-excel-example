package mongstub

import (
	"context"
	"os"
	"purchase-record/internal/database/mong"

	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// ContainerTerminateFunc defer this function when connect container for terminate container after used.
type ContainerTerminateFunc func(context.Context) error

// Container for stub mongo connection
type Container struct {
	Client    mong.IConnect
	Terminate ContainerTerminateFunc
}

// Connect is stub database connect, This database connect will be connect to mongo in container.
func Connect(databaseName string) (*Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:5.0.6",
		ExposedPorts: []string{"27017/tcp", "27018/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Waiting for connections"),
			wait.ForListeningPort("27017/tcp"),
		),
		Env: map[string]string{
			"MONGO_INITDB_ROOT_USERNAME": "root",
			"MONGO_INITDB_ROOT_PASSWORD": "example",
		},
	}

	// Creates a generic container with parameters
	container, err := testcontainers.GenericContainer(context.Background(), testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, errors.Wrap(err, "Creates a generic container with parameters")
	}

	// Get host where the container port is exposed
	host, err := container.Host(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "Get host where the container port is exposed")
	}

	// Get externally mapped port for a container port
	port, err := container.MappedPort(context.Background(), "27017/tcp")
	if err != nil {
		return nil, errors.Wrap(err, "Get externally mapped port for a container port")
	}

	// Make mongo connection string and set to environment variable
	os.Setenv("MONGO_URI", "mongodb://root:example@"+host+":"+port.Port()+"/")

	// New database connection
	client, err := mong.New(databaseName)
	if err != nil {
		return nil, errors.Wrap(err, "New database connection")
	}

	// Create a wrapper function for container termination
	terminateWrapper := func(ctx context.Context) error {
		return container.Terminate(ctx) // Call the original Terminate without additional options
	}

	// Done
	return &Container{
		Client:    client,
		Terminate: terminateWrapper,
	}, nil
}
