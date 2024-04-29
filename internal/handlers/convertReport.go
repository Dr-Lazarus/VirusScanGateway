package handlers

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Dr-Lazarus/VirusScanGateway/internal/db"
	vt "github.com/VirusTotal/vt-go"
)

func extractStringField(vtFile *vt.Object, field string) (string, error) {
	value, err := vtFile.GetString(field)
	if err != nil {
		log.Printf("Failed to get field '%s': %v. Using default empty string.", field, err)
		return "", nil
	}
	return value, nil
}

func extractTimeField(vtFile *vt.Object, field string) (time.Time, error) {
	t, err := vtFile.GetTime(field)
	if err != nil {
		log.Printf("Failed to get field '%s': %v. Using default zero time.", field, err)
		return time.Time{}, nil
	}
	return t, nil
}

func handleJSONField(data interface{}) json.RawMessage {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Printf("Failed to marshal JSON field: %v. Using default empty JSON.", err)
		return json.RawMessage("{}")
	}
	return json.RawMessage(bytes)
}

func ConvertToVirusTotalReport(vtFile *vt.Object) (*db.VirusTotalReport, error) {
	report := &db.VirusTotalReport{}

	var err error
	report.SHA256, _ = extractStringField(vtFile, "sha256")
	report.SHA1, _ = extractStringField(vtFile, "sha1")
	report.MD5, _ = extractStringField(vtFile, "md5")
	report.FileType, _ = extractStringField(vtFile, "type_tag")
	report.MeaningfulName, _ = extractStringField(vtFile, "meaningful_name")
	report.TypeDescription, _ = extractStringField(vtFile, "type_description")
	report.TypeExtension, _ = extractStringField(vtFile, "type_extension")

	report.Size, err = vtFile.GetInt64("size")
	if err != nil {
		log.Printf("Failed to get field 'size': %v. Using default zero.", err)
		report.Size = 0
	}

	report.FirstSubmissionDate, _ = extractTimeField(vtFile, "first_submission_date")
	report.LastSubmissionDate, _ = extractTimeField(vtFile, "last_submission_date")
	report.LastModificationDate, _ = extractTimeField(vtFile, "last_modification_date")
	report.LastAnalysisDate, _ = extractTimeField(vtFile, "last_analysis_date")

	if statsInterface, err := vtFile.Get("last_analysis_stats"); err == nil {
		report.LastAnalysisStats = handleJSONField(statsInterface)
	}
	if resultsInterface, err := vtFile.Get("last_analysis_results"); err == nil {
		report.LastAnalysisResults = handleJSONField(resultsInterface)
	}

	return report, nil
}
