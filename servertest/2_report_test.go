package servertest

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAndDeleteReportBySHA256(t *testing.T) {
	sha256Hashes := []string{
		"6254ae8c036a9d3c781f1b4a6f6940f192630888738a5942021e23e713f92e81",
		"f7f6a5894f1d19ddad6fa392b2ece2c5e578cbf7da4ea805b6885eb6985b6e3d",
		"73de4254959530e4d1d9bec586379184f96b4953dacf9cd5e5e2bdd7bfeceef7",
		"eaceb9628ee21df2e81494fdea8ee29e1c69e51777df13f35f824ab64d943ce4",
		"d824a0e03b01afb56b710f39f5713357613ccbd79eb2f7bedceaceecddb4c40e",
	}

	totalTests := 0
	passedTests := 0

	for _, hash := range sha256Hashes {
		t.Run("Get_SHA256_"+hash, func(t *testing.T) {
			if getReportBySHA256(t, hash) {
				passedTests++
			}
			totalTests++
		})
	}

	t.Logf("Completed %d get operations, %d passed", totalTests, passedTests)

	totalTests = 0
	passedTests = 0

	for _, hash := range sha256Hashes {
		t.Run("Delete_SHA256_"+hash, func(t *testing.T) {
			if deleteReportBySHA256(t, hash) {
				passedTests++
			}
			totalTests++
		})
	}
	t.Logf("Completed %d delete operations, %d passed", totalTests, passedTests)
}

func deleteReportBySHA256(t *testing.T, sha256 string) bool {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", "http://localhost:8080/reports/"+sha256, nil)
	if err != nil {
		t.Fatalf("Failed to create DELETE request: %s", err)
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send DELETE request: %s", err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
		return false
	}

	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		t.Fatalf("Failed to parse JSON response: %s", err)
		return false
	}

	if resp.StatusCode == http.StatusOK {
		assert.Equal(t, true, jsonResponse["success"], "Expected success to be true")
		assert.Equal(t, "Report deleted successfully.", jsonResponse["message"], "Expected successful deletion message")
		return true
	} else if resp.StatusCode == http.StatusNotFound {
		assert.Equal(t, false, jsonResponse["success"], "Expected success to be false")
		assert.Equal(t, "No report found with the given SHA256 ID", jsonResponse["error"].(map[string]interface{})["message"], "Expected no report found message")
		return false
	} else {
		t.Fatalf("Received unexpected HTTP status code: %d", resp.StatusCode)
		return false
	}
}

func getReportBySHA256(t *testing.T, sha256 string) bool {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/reports/"+sha256, nil)
	if err != nil {
		t.Fatalf("Failed to create GET request: %s", err)
		return false
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make GET request: %s", err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
		return false
	}

	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		t.Fatalf("Failed to parse JSON response: %s", err)
		return false
	}

	if resp.StatusCode == http.StatusOK {
		assert.Equal(t, true, jsonResponse["success"], "Expected success to be true")
		assert.Contains(t, []string{"Hold", "Completed"}, jsonResponse["processingStatus"], "Expected processing status to be 'Hold' or 'Completed'")
		assert.NotEmpty(t, jsonResponse["report"], "Expected report to be non-empty")
		return true
	} else if resp.StatusCode == http.StatusNotFound {
		assert.Equal(t, false, jsonResponse["success"], "Expected success to be false")
		assert.Equal(t, "No report found with the given SHA256 ID", jsonResponse["error"].(map[string]interface{})["message"], "Expected no report found message")
		return false
	} else if resp.StatusCode == http.StatusInternalServerError {
		assert.Equal(t, false, jsonResponse["success"], "Expected success to be false")
		assert.Equal(t, "Error retrieving report from the database", jsonResponse["error"].(map[string]interface{})["message"], "Expected internal server error message")
		return false
	} else {
		t.Fatalf("Received unexpected HTTP status code: %d", resp.StatusCode)
		return false
	}
}
