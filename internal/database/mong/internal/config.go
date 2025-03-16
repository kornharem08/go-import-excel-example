package internal

type Config struct {
	URI string `env:"MONGO_URI" default:"mongodb://localhost:27017"`
}
