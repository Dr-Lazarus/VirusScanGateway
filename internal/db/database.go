package db

import (
	"database/sql"
	"log"

	"github.com/Dr-Lazarus/VirusScanGateway/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
)

func SetupDatabase(cfg *config.Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging the database: ", err)
	}

	log.Println("✅ Successfully connected to the database.")
	runMigrations(cfg.DatabaseURL)

	return db
}

func runMigrations(connStr string) {
	m, err := migrate.New("file://internal/db/migrations", connStr)
	if err != nil {
		log.Fatal("Error creating migration: ", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Error applying migration: ", err)
	}

	log.Println("🎉 Database migrated successfully.")
}
