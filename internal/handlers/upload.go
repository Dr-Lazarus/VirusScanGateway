package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Dr-Lazarus/VirusScanGateway/pkg/api/virustotal" // Adjust the path as needed
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func uploadHandler(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.String(http.StatusMethodNotAllowed, "Unsupported method")
		return
	}

	if os.Getenv("APP_ENV") == "DEV" {
		err := godotenv.Load(".env.dev")
		if err != nil {
			c.String(http.StatusInternalServerError, "Error loading .env.dev file")
			return
		}
	}

	apiKey := os.Getenv("VIRUSTOTAL_API_KEY")
	if apiKey == "" {
		c.String(http.StatusInternalServerError, "VirusTotal API key not set")
		return
	}

	client := virustotal.NewClient(apiKey)

	const maxUploadSize = 10 * 1024 * 1024
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid file")
		return
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading file")
		return
	}

	response, err := client.UploadFile(header.Filename, fileContent)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error uploading file to VirusTotal: %s", err))
		return
	}

	c.String(http.StatusOK, "VirusTotal Response: %s", response)
}
