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
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmplPath := filepath.Join("web", "templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "⚠️ Internal Server Error", 500)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "⚠️ Internal Server Error", 500)
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

	if env == "DEV" {
		if err := runMigrations(connStr); err != nil {
			log.Fatal("❌ Error running migrations: ", err)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("❌ $PORT not set")
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/upload", handler.UploadHandler)
	log.Printf("🚀 Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("❌ Server failed to start: ", err)
	}
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
