package handlers

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"

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

	const maxUploadSize = 10 * 1024 * 1024
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)
	err := c.Request.ParseMultipartForm(maxUploadSize)
	if err != nil {
		c.String(http.StatusBadRequest, "The uploaded file is too big. Please choose a file that's less than 10MB in size")
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid file")
		return
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error processing file")
		return
	}
	_, err = io.Copy(part, file)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error processing file")
		return
	}
	writer.Close()

	req, err := http.NewRequest("POST", "https://www.virustotal.com/api/v3/files", body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error creating request to VirusTotal API")
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-apikey", apiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error sending file to VirusTotal API")
		return
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading response from VirusTotal API")
		return
	}

	c.String(http.StatusOK, "VirusTotal Response: %s", responseBody)
}

func RegisterRoutes(router *gin.Engine) {
	router.GET("/upload", uploadHandler)
}
