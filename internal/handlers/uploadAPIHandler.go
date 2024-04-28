package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

func uploadVirusHandler(file io.Reader, filename, apiKey string) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("error creating form file: %v", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return "", fmt.Errorf("error copying file data: %v", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", "https://www.virustotal.com/api/v3/files", body)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-apikey", apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send the request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %v", err)
	}
	if analysisURL, ok := result["data"].(map[string]interface{})["links"].(map[string]interface{})["self"].(string); ok {
		return analysisURL, nil
	}
	return "", fmt.Errorf("failed to extract analysis URL")
}
