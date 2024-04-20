package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB) {
	router.GET("/", homeHandler)
	router.POST("/upload", func(c *gin.Context) {
		uploadHandler(c, db)
		// uploadHandler(c)
	})
	router.GET("/getbyID/:SHA256ID", func(c *gin.Context) {
		getidHandler(c, db)
	})
}
