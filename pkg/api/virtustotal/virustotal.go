package virustotal

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

// UploadFile uploads a file to the VirusTotal API and returns the API response
func (c *Client) UploadFile(filename string, fileContent []byte) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(part, bytes.NewReader(fileContent))
	if err != nil {
		return "", err
	}
	writer.Close()

	headers := map[string]string{
		"Content-Type": writer.FormDataContentType(),
	}
	response, err := c.makeRequest(http.MethodPost, "https://www.virustotal.com/api/v3/files", body, headers)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(respBody), nil
}
