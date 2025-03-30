package main

import (
	_ "purchase-record/docs"
	"purchase-record/internal/router"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Excel Order API
// @version 1.0
// @description API for processing Excel order data
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@example.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @schemes http https
// @in header
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	r := gin.Default()

	// Configure CORS
	r.Use(cors.Default())

	// Routes
	router.RegisterRoutePurchaseOrder(r)

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run server
	r.Run(":8080")
}
