package handlers

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/gin-gonic/gin"
)

func homeHandler(c *gin.Context) {
	tmplPath := filepath.Join("web", "templates", "index.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		c.String(http.StatusInternalServerError, "❌ Internal Server Error")
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err = tmpl.Execute(c.Writer, nil); err != nil {
		c.String(http.StatusInternalServerError, "❌ Internal Server Error")
	}
}
