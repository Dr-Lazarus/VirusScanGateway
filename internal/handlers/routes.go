package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/", homeHandler)
	router.POST("/upload", func(c *gin.Context) {
		uploadHandler(c, db)
	})
	router.GET("/reports/:SHA256ID", func(c *gin.Context) {
		getReportHandler(c, db)
	})
	router.DELETE("/reports/:SHA256ID", func(c *gin.Context) {
		deleteReportHandler(c, db)
	})
}
