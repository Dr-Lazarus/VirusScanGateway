package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/Dr-Lazarus/VirusScanGateway/internal/db"
	"github.com/gin-gonic/gin"
)

func homeHandler(c *gin.Context) {
	tmplPath := filepath.Join("web", "templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "⚠️ Internal Server Error")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err = tmpl.Execute(c.Writer, nil); err != nil {
		c.String(http.StatusInternalServerError, "⚠️ Internal Server Error")
	}
}

func getidHandler(c *gin.Context, dbConn *sql.DB) {
	// Retrieve SHA256 from the URL path parameter
	sha256 := c.Param("SHA256ID")
	log.Println("SHA256:", sha256)

	// Use the retrieved SHA256 to get the report from the database
	report, err := db.GetReport(dbConn, sha256)
	if err != nil {
		if err == sql.ErrNoRows {
			c.String(http.StatusNotFound, "No report found with the given SHA256 ID")
		} else {
			log.Printf("❌ Error retrieving report: %v", err)
			c.String(http.StatusInternalServerError, "❌ Error retrieving report from the database")
		}
		return
	}

	// Check if the LastAnalysisResults are empty and inform the user accordingly
	if len(report.LastAnalysisResults) == 0 || string(report.LastAnalysisResults) == "{}" {
		c.String(http.StatusOK, "Report has not yet been updated. If 3 minutes have passed since you uploaded the report, please reupload the report.")
	} else {
		c.JSON(http.StatusOK, report)
	}
}
