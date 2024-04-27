package main

import (
	"log"
	"net/http"
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

	env := os.Getenv("ENVIRONMENT")

	if env == "PROD" {
		httpsPort := os.Getenv("HTTPS_PORT")
		if httpsPort == "" {
			httpsPort = "443"
		}

		sslCert := os.Getenv("SSL_CERT_PATH")
		sslKey := os.Getenv("SSL_KEY_PATH")
		log.Println("[DEBUG]")

		if sslCert == "" || sslKey == "" {
			log.Fatal("❌ SSL certificate or key file path not set")
		}

		go func() {
			httpPort := os.Getenv("HTTP_PORT")
			if httpPort == "" {
				httpPort = "80"
			}

			log.Printf("🚀 HTTP server starting on port %s for redirect to HTTPS\n", httpPort)
			httpRouter := gin.Default()
			// Redirect all HTTP requests to HTTPS
			httpRouter.GET("/*any", func(c *gin.Context) {
				target := "https://" + c.Request.Host + c.Request.URL.Path
				if len(c.Request.URL.RawQuery) > 0 {
					target += "?" + c.Request.URL.RawQuery
				}
				c.Redirect(http.StatusMovedPermanently, target)
			})
			log.Fatal(http.ListenAndServe(":"+httpPort, httpRouter))
		}()

		log.Printf("🚀 HTTPS server starting on port %s\n", httpsPort)
		log.Fatal(http.ListenAndServeTLS(":"+httpsPort, sslCert, sslKey, router))
	} else {
		httpPort := os.Getenv("HTTP_PORT")
		if httpPort == "" {
			httpPort = "8080"
		}

		log.Printf("🚀 HTTP server starting on port %s in development mode\n", httpPort)
		log.Fatal(http.ListenAndServe(":"+httpPort, router))
	}
}
