package main

import (
	"log"
	_ "purchase-record/docs"
	"purchase-record/internal/config"
	"purchase-record/internal/database/environ"
	"purchase-record/internal/database/mong"
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
func main() {
	cfg := environ.Load[config.Config]()
	if cfg.MongoDBDatabase == "" {
		log.Fatal("database name must be present")
	}

	dbconn, err := mong.New(cfg.MongoDBDatabase)
	// New database connection
	if err != nil {
		log.Fatal(err)
	}

	// Ensure connection close
	defer dbconn.Close()

	r := gin.Default()
	r.Use(cors.Default())

	// Routes
	router.RegisterRoutePurchaseOrder(r, dbconn)
	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run server
	r.Run(":8080")
}
