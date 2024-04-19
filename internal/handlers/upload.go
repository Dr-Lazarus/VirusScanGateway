package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Dr-Lazarus/VirusScanGateway/internal/db"
	vt "github.com/VirusTotal/vt-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func uploadHandler(c *gin.Context, dbConn *sql.DB) {
	log.Println("Handling file upload...")

	if c.Request.Method != http.MethodPost {
		log.Println("Error: Method not allowed")
		c.String(http.StatusMethodNotAllowed, "Unsupported method")
		return
	}

	if os.Getenv("APP_ENV") == "DEV" {
		if err := godotenv.Load(".env.dev"); err != nil {
			log.Printf("Error loading .env.dev file: %v", err)
			c.String(http.StatusInternalServerError, "Error loading .env.dev file")
			return
		}
	}

	apiKey := os.Getenv("VIRUSTOTAL_API_KEY")
	if apiKey == "" {
		log.Println("Error: API key is not set")
		c.String(http.StatusInternalServerError, "VirusTotal API key not set")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Println("Error retrieving the file")
		c.String(http.StatusBadRequest, "Invalid file upload")
		return
	}
	defer file.Close()

	analysisURL, err := uploadVirusHandler(file, header.Filename, apiKey)
	if err != nil {
		log.Printf("Error uploading to VirusTotal: %v", err)
		c.String(http.StatusInternalServerError, "Failed to upload file to VirusTotal")
		return
	}

	sha256, err := analysisHandler(analysisURL, apiKey)
	if err != nil {
		log.Printf("Error extracting SHA256: %v", err)
		c.String(http.StatusInternalServerError, "Failed to extract SHA256")
		return
	}

	client := vt.NewClient(apiKey)
	vtFile, err := client.GetObject(vt.URL("files/%s", sha256))
	if err != nil {
		log.Printf("Failed to retrieve file details: %v", err)
		c.String(http.StatusInternalServerError, "Failed to retrieve file details")
		return
	}
	report, err := ConvertToVirusTotalReport(vtFile)
	if err != nil {
		log.Printf("Failed to convert virus total report: %v", err)
		c.String(http.StatusInternalServerError, "Failed to process file details")
		return
	}

	_, err = db.GetReport(dbConn, report.SHA256)
	if err != nil {
		if err == sql.ErrNoRows {
			err = db.InsertReport(dbConn, report)
			if err != nil {
				log.Printf("Failed to insert report: %v", err)
				c.String(http.StatusInternalServerError, "Failed to insert report")
				return
			}
		} else {
			log.Printf("Error checking for existing report: %v", err)
			c.String(http.StatusInternalServerError, "Error checking for existing report")
			return
		}
	} else {
		err = db.UpdateReport(dbConn, report)
		if err != nil {
			log.Printf("Failed to update report: %v", err)
			c.String(http.StatusInternalServerError, "Failed to update report")
			return
		}
	}

	c.JSON(http.StatusOK, report)

}
