package main

import (
	"log"
	"os"

	"github.com/Dr-Lazarus/VirusScanGateway/internal/config"
	"github.com/Dr-Lazarus/VirusScanGateway/internal/db"
	"github.com/Dr-Lazarus/VirusScanGateway/internal/handlers"
	"github.com/Dr-Lazarus/VirusScanGateway/internal/middleware"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	router := gin.Default()
	router.Use(middleware.CORSMiddleware())

	dbConn := db.SetupDatabase(cfg)
	defer dbConn.Close()

	handlers.RegisterRoutes(router, dbConn)
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("❌ $PORT not set")
	}

	log.Printf("🚀 Server starting on port %s\n", port)
	router.Run(":" + port)
}
