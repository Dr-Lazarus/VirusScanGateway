package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func makeGetRequest(analysisURL, apiKey string) ([]byte, error) {
	req, err := http.NewRequest("GET", analysisURL, nil)
	if err != nil {
		return nil, fmt.Errorf("❌ Error creating request: %v", err)
	}

	req.Header.Set("x-apikey", apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to send the request: %v", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("❌ Failed to read response: %v", err)
	}

	return respBody, nil
}

func analysisHandler(analysisURL, apiKey string) (string, error) {
	respBody, err := makeGetRequest(analysisURL, apiKey)
	if err != nil {
		return "", fmt.Errorf("❌ Failed to get analysis information: %v", err)
	}

	var result AnalysisResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("❌ Failed to parse response: %v", err)
	}

	sha256 := result.Meta.FileInfo.SHA256
	log.Println("[INFO] SHA256:", result.Meta.FileInfo.SHA256)

	return sha256, nil
}
