package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Dr-Lazarus/VirusScanGateway/internal/handler"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Ensure this line is exactly as shown
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func homeHandler(c *gin.Context) {
	tmplPath := filepath.Join("web", "templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "⚠️ Internal Server Error")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(c.Writer, nil)
	if err != nil {
		c.String(http.StatusInternalServerError, "⚠️ Internal Server Error")
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	env := os.Getenv("APP_ENV")
	if env == "DEV" {
		if err := godotenv.Load(".env.dev"); err != nil {
			log.Fatal("❌ Error loading .env.dev file")
		}
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Error connecting to the database: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("❌ Error pinging the database: ", err)
	} else {
		log.Println("✅ Successfully connected to the database.")
	}

	if err := runMigrations(connStr); err != nil {
		log.Fatal("❌ Error running migrations: ", err)
	}

	router := gin.Default()
	router.Use(CORSMiddleware())
	router.GET("/", homeHandler)
	router.POST("/upload", handler.UploadHandler)
	router.Static("/static", "./web/static")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("❌ $PORT not set")
	}

	log.Printf("🚀 Server starting on port %s\n", port)
	router.Run(":" + port) // Use Gin to run the server on the correct port
}

func runMigrations(connStr string) error {
	m, err := migrate.New(
		"file://pkg/database/migrations",
		connStr,
	)
	if err != nil {
		return fmt.Errorf("❌ Error creating migration: %w", err)
	}

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("🆗 Table up to date with latest schema")
			return nil
		}
		return fmt.Errorf("❌ Error applying migration: %w", err)
	}

	log.Println("🎉 Database migrated successfully.")
	return nil
}
