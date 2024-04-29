package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Dr-Lazarus/VirusScanGateway/internal/db"
	vt "github.com/VirusTotal/vt-go"
	"github.com/gin-gonic/gin"
)

func uploadHandler(c *gin.Context, dbConn *sql.DB) {
	log.Println("Handling file upload...")

	if c.Request.Method != http.MethodPost {
		log.Println("❌ Error: Method not allowed")
		c.String(http.StatusMethodNotAllowed, "Unsupported method")
		return
	}

	apiKey := os.Getenv("VIRUSTOTAL_API_KEY")
	if apiKey == "" {
		log.Println("❌ Error: API key is not set")
		c.String(http.StatusInternalServerError, "VirusTotal API key not set")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		log.Println("❌ Error retrieving the file")
		c.String(http.StatusBadRequest, "Invalid file upload")
		return
	}
	defer file.Close()

	analysisURL, err := uploadVirusHandler(file, header.Filename, apiKey)
	if err != nil {
		log.Printf("❌ Error uploading to VirusTotal: %v", err)
		c.String(http.StatusInternalServerError, "Failed to upload file to VirusTotal")
		return
	}

	sha256, err := analysisHandler(analysisURL, apiKey)
	if err != nil {
		log.Printf("❌ Error extracting SHA256: %v", err)
		c.String(http.StatusInternalServerError, "Failed to extract SHA256")
		return
	}
	reportSQL, err := db.GetReport(dbConn, sha256)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("❌ Error checking for existing report: %v", err)
		c.String(http.StatusInternalServerError, "❌ Error checking for existing report")
		return
	}

	if err == nil {
		if len(reportSQL.LastAnalysisResults) != 0 && string(reportSQL.LastAnalysisResults) != "{}" {
			message := "Report already present in DB. Please use SHA256 ID to retreive report."
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": message,
				"sha256":  sha256,
			})
			return
		} else if (len(reportSQL.LastAnalysisResults) == 0 || string(reportSQL.LastAnalysisResults) == "{}") && IsProcessing(sha256) {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Please Hold 🫷. Your file is still being processed.",
				"sha256":  sha256,
			})
			return
		}
	}

	client := vt.NewClient(apiKey)
	vtFile, err := client.GetObject(vt.URL("files/%s", sha256))
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.String(http.StatusNotFound, "❌ Error Uploading File. Please upload the file again.")
		} else {
			c.String(http.StatusInternalServerError, "An error occurred while processing the file. Please try uploading the file again.")
		}
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
				log.Printf("❌ Failed to insert report: %v", err)
				c.String(http.StatusInternalServerError, "Failed to insert report")
				return
			}
		} else {
			log.Printf("❌ Error checking for existing report: %v", err)
			c.String(http.StatusInternalServerError, "Error checking for existing report")
			return
		}
	} else {
		err = db.UpdateReport(dbConn, report)
		if err != nil {
			log.Printf("❌ Failed to update report: %v", err)
			c.String(http.StatusInternalServerError, "Failed to update report")
			return
		}
	}
	if len(report.LastAnalysisResults) == 0 || string(report.LastAnalysisResults) == "{}" {
		SetProcessing(sha256, true)
		go backgroundWorker(dbConn, apiKey, sha256)
	}

	message := fmt.Sprintf("The file has been uploaded and is now being processed by the Virus Total API." +
		"Please allow some time for the analysis to complete." +
		"You can check the status of the report using the provided SHA256 ID.")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"sha256":  sha256,
	})

}

func backgroundWorker(dbConn *sql.DB, apiKey, sha256 string) {
	defer ClearProcessing(sha256)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	client := vt.NewClient(apiKey)
	attemptCounter := 0

	for {
		select {
		case <-ticker.C:
			vtFile, err := client.GetObject(vt.URL("files/%s", sha256))
			if err != nil {
				log.Printf("❌ Failed to retrieve file details: %v", err)
				continue
			}
			report, err := ConvertToVirusTotalReport(vtFile)
			if err != nil {
				log.Printf("❌ Failed to convert virus total report: %v", err)
				continue
			}

			if len(report.LastAnalysisResults) == 0 || string(report.LastAnalysisResults) == "{}" {
				log.Println("Last analysis results are empty, continuing to poll.")
				attemptCounter++
				if attemptCounter >= 2 {
					log.Printf("Stopping after 2 attempts.")
					return
				}
				continue
			}

			err = db.UpdateReport(dbConn, report)
			if err != nil {
				log.Printf("❌ Failed to update report: %v", err)
				continue
			}

			log.Println("Report updated successfully. SHA 256:", report.SHA256)
			return
		}
	}
}
