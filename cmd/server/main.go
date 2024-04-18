package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"

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
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
	}
}

func main() {
	env := os.Getenv("APP_ENV")
	var envFile string
	if env == "PROD" {
		envFile = ".env.prod"
	} else {
		envFile = ".env.dev"
	}

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("Error loading environment file: ", envFile)
	}

	var connStr string
	if env == "PROD" {
		connStr = fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=require",
			os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
	} else {
		connStr = fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable host=localhost",
			os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	} else {
		log.Println("Successfully connected to the database.")
	}

	http.HandleFunc("/", homeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("$PORT not set, defaulting to %s in development", port)
	}

	log.Printf("Server starting on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
