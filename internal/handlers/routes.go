package handlers

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/", homeHandler)
	router.POST("/upload", uploadHandler)
}
