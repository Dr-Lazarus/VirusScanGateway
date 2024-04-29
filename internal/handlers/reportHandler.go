package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/Dr-Lazarus/VirusScanGateway/internal/db"
	"github.com/gin-gonic/gin"
)

func getReportHandler(c *gin.Context, dbConn *sql.DB) {
	sha256 := c.Param("SHA256ID")
	log.Println("SHA256:", sha256)

	report, err := db.GetReport(dbConn, sha256)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error": gin.H{
					"message": "No report found with the given SHA256 ID",
					"code":    http.StatusNotFound,
				},
			})
		} else {
			log.Printf("❌ Error retrieving report: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"message": "Error retrieving report from the database",
					"code":    http.StatusInternalServerError,
				},
			})
		}
		return
	}

	var processingStatus string
	if len(report.LastAnalysisResults) == 0 || string(report.LastAnalysisResults) == "{}" {
		processingStatus = "Hold"
	} else {
		processingStatus = "Completed"
	}

	c.JSON(http.StatusOK, gin.H{
		"success":          true,
		"processingStatus": processingStatus,
		"report":           report,
	})
}

func deleteReportHandler(c *gin.Context, dbConn *sql.DB) {
	sha256 := c.Param("SHA256ID")
	log.Println("SHA256:", sha256)

	_, err := db.GetReport(dbConn, sha256)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error": gin.H{
					"message": "No report found with the given SHA256 ID",
					"code":    http.StatusNotFound,
				},
			})
		} else {
			log.Printf("❌ Error retrieving report: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error": gin.H{
					"message": "Error retrieving report from the database",
					"code":    http.StatusInternalServerError,
				},
			})
		}
		return
	}
	err = db.DeleteReport(dbConn, sha256)
	if err != nil {
		log.Printf("❌ Error deleting report: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error": gin.H{
				"message": "Error deleting report from the database",
				"code":    http.StatusInternalServerError,
			},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Report deleted successfully.",
	})
}
