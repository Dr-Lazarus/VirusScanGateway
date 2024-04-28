package servertest

import (
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHomePage(t *testing.T) {
	// Initialize a counter for passed checks
	passedChecks := 0
	totalChecks := 2 // Update this if you add more checks

	resp, err := http.Get("http://127.0.0.1:8080/")
	if err != nil {
		t.Fatalf("Failed to make GET request to the home page: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	} else {
		passedChecks++
		t.Log("Passed: Status code check")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %s", err)
	}

	expected, err := os.ReadFile("/home/runner/work/VirusScanGateway/VirusScanGateway/web/templates/index.html")
	if err != nil {
		t.Fatalf("Failed to read expected file: %s", err)
	}

	if assert.Equal(t, string(expected), string(body), "The HTML content does not match the expected content.") {
		passedChecks++
		t.Log("Passed: Content check")
	}

	t.Logf("Test cases passed: %d/%d", passedChecks, totalChecks)
	if passedChecks != totalChecks {
		t.Errorf("Some test cases failed: %d/%d", totalChecks-passedChecks, totalChecks)
	}
}
