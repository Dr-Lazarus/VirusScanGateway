package servertest

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadMultipleFiles(t *testing.T) {
	dirPath := "../test_files"

	files, err := os.ReadDir(dirPath)
	if err != nil {
		t.Fatalf("Failed to read directory: %s", err)
	}

	expectedMessages := map[string]bool{
		"Report already present in DB. Please use SHA256 ID to retreive report.": true,
		"The file has been uploaded and is now being processed by the Virus Total API. Please allow some time for the analysis to complete. You can check the status of the report using the provided SHA256 ID.": true,
		"Please Hold 🫷. Your file is still being processed.": true,
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		t.Run(file.Name(), func(t *testing.T) {
			filePath := filepath.Join(dirPath, file.Name())

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, err := writer.CreateFormFile("file", file.Name())
			if err != nil {
				t.Fatalf("Failed to create form file: %s", err)
			}

			fileData, err := os.ReadFile(filePath)
			if err != nil {
				t.Fatalf("Failed to read test file: %s", err)
			}
			part.Write(fileData)
			writer.Close()

			req, err := http.NewRequest("POST", "http://127.0.0.1:8080/upload", body)
			if err != nil {
				t.Fatalf("Failed to create POST request: %s", err)
			}
			req.Header.Set("Content-Type", writer.FormDataContentType())

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Failed to send POST request: %s", err)
			}
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected HTTP status 200 OK for file: "+file.Name())

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %s", err)
			}

			var jsonResponse map[string]interface{}
			if err := json.Unmarshal(respBody, &jsonResponse); err != nil {
				t.Fatalf("Failed to parse JSON response: %s", err)
			}

			assert.True(t, jsonResponse["success"].(bool), "Expected success to be true for file: "+file.Name())
			assert.NotEmpty(t, jsonResponse["sha256"], "Expected a non-empty SHA256 ID for file: "+file.Name())

			message := jsonResponse["message"].(string)
			log.Println("[DEBUG] Message:", message)
			_, ok := expectedMessages[message]
			log.Println("[DEBUG] Expected Message:", expectedMessages[message])
			assert.True(t, ok, "Received unexpected message for file: "+file.Name()+" - "+message)
		})
	}
}
