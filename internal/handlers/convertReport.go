package handlers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Dr-Lazarus/VirusScanGateway/internal/db"
	vt "github.com/VirusTotal/vt-go"
)

func ConvertToVirusTotalReport(vtFile *vt.Object) (*db.VirusTotalReport, error) {
	report := &db.VirusTotalReport{}

	// Attempt to extract fields using direct Get calls and handling the returned error
	extractStringField := func(field string) (string, error) {
		value, err := vtFile.GetString(field)
		if err != nil {
			log.Printf("Failed to get field '%s': %v. Using default empty string.", field, err)
			return "", nil // Return empty string if there is an error
		}
		return value, nil
	}

	var err error
	report.SHA256, _ = extractStringField("sha256")
	report.SHA1, _ = extractStringField("sha1")
	report.MD5, _ = extractStringField("md5")
	report.FileType, _ = extractStringField("type_tag")
	report.MeaningfulName, _ = extractStringField("meaningful_name")
	report.TypeDescription, _ = extractStringField("type_description")
	report.TypeExtension, _ = extractStringField("type_extension")

	report.Size, err = vtFile.GetInt64("size")
	if err != nil {
		log.Printf("Failed to get field 'size': %v. Using default zero.", err)
		report.Size = 0 // Use default zero value if there is an error
	}

	// Helper function to handle time fields
	extractTimeField := func(field string) (time.Time, error) {
		t, err := vtFile.GetTime(field)
		if err != nil {
			log.Printf("Failed to get field '%s': %v. Using default zero time.", field, err)
			return time.Time{}, nil // Return zero time if there is an error
		}
		return t, nil
	}

	report.FirstSubmissionDate, _ = extractTimeField("first_submission_date")
	report.LastSubmissionDate, _ = extractTimeField("last_submission_date")
	report.LastModificationDate, _ = extractTimeField("last_modification_date")
	report.LastAnalysisDate, _ = extractTimeField("last_analysis_date")

	// JSONB fields handling
	handleJSONField := func(data interface{}) json.RawMessage {
		if bytes, err := json.Marshal(data); err == nil {
			return json.RawMessage(bytes)
		}
		log.Printf("Failed to marshal JSON field: %v. Using default empty JSON.", err)
		return json.RawMessage("{}") // Return empty JSON if there is an error
	}

	if statsInterface, err := vtFile.Get("last_analysis_stats"); err == nil {
		report.LastAnalysisStats = handleJSONField(statsInterface)
	}

	if resultsInterface, err := vtFile.Get("last_analysis_results"); err == nil {
		report.LastAnalysisResults = handleJSONField(resultsInterface)
	}

	return report, nil
}
