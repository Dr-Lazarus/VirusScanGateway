package db

import (
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/Dr-Lazarus/VirusScanGateway/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type VirusTotalReport struct {
	FileType             string          `json:"file_type"`
	SHA256               string          `json:"sha256"`
	SHA1                 string          `json:"sha1"`
	MD5                  string          `json:"md5"`
	MeaningfulName       string          `json:"meaningful_name"`
	TypeDescription      string          `json:"type_description"`
	TypeExtension        string          `json:"type_extension"`
	LastModificationDate time.Time       `json:"last_modification_date"`
	FirstSubmissionDate  time.Time       `json:"first_submission_date"`
	LastSubmissionDate   time.Time       `json:"last_submission_date"`
	Size                 int64           `json:"size"`
	LastAnalysisDate     time.Time       `json:"last_analysis_date"`
	LastAnalysisStats    json.RawMessage `json:"last_analysis_stats"`
	LastAnalysisResults  json.RawMessage `json:"last_analysis_results"`
}

func SetupDatabase(cfg *config.Config) *sql.DB {
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("❌ Error connecting to the database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Error pinging the database: %v", err)
	}

	log.Println("✅ Successfully connected to the database.")
	runMigrations(cfg.DatabaseURL)

	return db
}

func runMigrations(connStr string) {
	m, err := migrate.New("file://pkg/database/migrations", connStr)
	if err != nil {
		log.Fatalf("❌ Error creating migration: %v", err)
	}
	err = m.Up()
	if err == migrate.ErrNoChange {
		log.Println("✅ Database synced up")
	} else if err != nil {
		log.Fatalf("❌ Error applying migration: %v", err)
	} else {
		log.Println("🎉 Database migrated successfully.")
	}
}

func InsertReport(db *sql.DB, report *VirusTotalReport) error {
	query := `
        INSERT INTO virustotal_reports (
            file_type, sha256, sha1, md5, meaningful_name, type_description, type_extension,
            last_modification_date, first_submission_date, last_submission_date, size,
            last_analysis_date, last_analysis_stats, last_analysis_results
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
        ON CONFLICT (sha256) DO UPDATE SET
            file_type = EXCLUDED.file_type, sha1 = EXCLUDED.sha1, md5 = EXCLUDED.md5,
            meaningful_name = EXCLUDED.meaningful_name, type_description = EXCLUDED.type_description,
            type_extension = EXCLUDED.type_extension, last_modification_date = EXCLUDED.last_modification_date,
            first_submission_date = EXCLUDED.first_submission_date, last_submission_date = EXCLUDED.last_submission_date,
            size = EXCLUDED.size, last_analysis_date = EXCLUDED.last_analysis_date,
            last_analysis_stats = EXCLUDED.last_analysis_stats, last_analysis_results = EXCLUDED.last_analysis_results;
    `
	_, err := db.Exec(query, report.FileType, report.SHA256, report.SHA1, report.MD5, report.MeaningfulName,
		report.TypeDescription, report.TypeExtension, report.LastModificationDate, report.FirstSubmissionDate,
		report.LastSubmissionDate, report.Size, report.LastAnalysisDate,
		report.LastAnalysisStats, report.LastAnalysisResults)
	return err
}

func GetReport(db *sql.DB, sha256 string) (*VirusTotalReport, error) {
	var report VirusTotalReport
	query := `SELECT * FROM virustotal_reports WHERE sha256 = $1;`
	err := db.QueryRow(query, sha256).Scan(
		&report.FileType, &report.SHA256, &report.SHA1, &report.MD5,
		&report.MeaningfulName, &report.TypeDescription, &report.TypeExtension,
		&report.LastModificationDate, &report.FirstSubmissionDate, &report.LastSubmissionDate,
		&report.Size, &report.LastAnalysisDate,
		&report.LastAnalysisStats, &report.LastAnalysisResults,
	)
	if err != nil {
		return nil, err
	}
	return &report, nil
}

func UpdateReport(db *sql.DB, report *VirusTotalReport) error {
	query := `
        UPDATE virustotal_reports SET
            file_type = $1, sha1 = $2, md5 = $3, meaningful_name = $4,
            type_description = $5, type_extension = $6, last_modification_date = $7,
            first_submission_date = $8, last_submission_date = $9, size = $10,
            last_analysis_date = $11, last_analysis_stats = $12, last_analysis_results = $13
        WHERE sha256 = $14;
    `
	_, err := db.Exec(query, report.FileType, report.SHA1, report.MD5, report.MeaningfulName,
		report.TypeDescription, report.TypeExtension, report.LastModificationDate, report.FirstSubmissionDate,
		report.LastSubmissionDate, report.Size, report.LastAnalysisDate,
		report.LastAnalysisStats, report.LastAnalysisResults, report.SHA256)
	return err
}

func DeleteReport(db *sql.DB, sha256 string) error {
	query := `DELETE FROM virustotal_reports WHERE sha256 = $1;`
	_, err := db.Exec(query, sha256)
	return err
}
