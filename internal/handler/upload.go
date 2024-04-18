package handler

import (
	"fmt"
	"io"
	"net/http"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "❌ Unsupported method", http.StatusMethodNotAllowed)
		return
	}

	const maxUploadSize = 10 * 1024 * 1024
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "❌ The uploaded file is too big. Please choose an file that's less than 10MB in size", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "❌ Invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "❌ Invalid file", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "🎉 File contents:\n\n %s", fileBytes)
}
