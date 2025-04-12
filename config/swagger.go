package config

// SwaggerConfig contains configuration for Swagger documentation
type SwaggerConfig struct {
	Title       string
	Description string
	Version     string
	Host        string
	BaseURL     string
}

// Config holds all application configurations
type Config struct {
	Swagger SwaggerConfig
}

// CF is the global configuration instance
var CF Config

// InitSwaggerConfig initializes swagger configuration
func InitSwaggerConfig() {
	CF.Swagger = SwaggerConfig{
		Title:       "Excel Order API",
		Description: "API for processing Excel order data",
		Version:     "1.0",
		Host:        "localhost:8080", //localhost:8080 or 10.10.50.5:8080,
		BaseURL:     "",
	}
}
