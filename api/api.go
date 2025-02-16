package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type MediaEntry struct {
	SKU  string `json:"sku"`
	Path string `json:"path"`
}

func GetPaths() []MediaEntry {
	mediaHandlerURL := "https://api.digiapi.org/apiv2/mediaHandler"
	resp, err := http.Get(mediaHandlerURL)
	if err != nil {
		log.Fatalf("Error fetching from %s: %v", mediaHandlerURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response: %v", err)
	}

	var data []MediaEntry
	if err := json.Unmarshal(body, &data); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	return data
}

func UploadImage(data MediaEntry) error {
	url := "https://api.digiapi.org/apiv2/mediaUploader"
	file, err := os.Open(data.Path)
	if err != nil {
		return fmt.Errorf("failed to open file %q: %w", data.Path, err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	err = writer.WriteField("sku", data.SKU)
	if err != nil {
		return fmt.Errorf("failed to write sku field: %w", err)
	}

	part, err := writer.CreateFormFile("image", file.Name())
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload failed: status %s, response: %s", resp.Status, respBody)
	}

	return nil
}
