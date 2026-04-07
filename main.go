package main

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FileType int

const (
	TypeDocument FileType = iota
	TypePhoto
	TypeVideo
)

func detectFileType(file *os.File) (FileType, string) {
	buffer := make([]byte, 512)
	n, _ := file.ReadAt(buffer, 0)
	mimeType := http.DetectContentType(buffer[:n])

	if strings.HasPrefix(mimeType, "image/") {
		return TypePhoto, "photo"
	}
	if strings.HasPrefix(mimeType, "video/") {
		return TypeVideo, "video"
	}

	ext := strings.ToLower(filepath.Ext(file.Name()))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return TypePhoto, "photo"
	case ".mp4", ".mov", ".m4v", ".avi", ".mkv":
		return TypeVideo, "video"
	default:
		return TypeDocument, "document"
	}
}

func getEndpoint(ft FileType) string {
	switch ft {
	case TypePhoto:
		return "sendPhoto"
	case TypeVideo:
		return "sendVideo"
	default:
		return "sendDocument"
	}
}

func sendFile(cfg *Config, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	fileType, fieldName := detectFileType(file)
	endpoint := getEndpoint(fileType)
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", cfg.Token, endpoint)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("chat_id", cfg.ChatID); err != nil {
		return err
	}

	part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
	if err != nil {
		return err
	}
	if _, err := io.Copy(part, file); err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  stt <file1> [file2] ...  Send files to Telegram")
		fmt.Println("  stt --setup              Configure bot token and chat ID")
		os.Exit(1)
	}

	if os.Args[1] == "--setup" {
		if err := setup(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error during setup: %v\n", err)
			os.Exit(1)
		}
		return
	}

	cfg, err := loadConfig()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v. Run 'stt --setup' to configure.\n", err)
		os.Exit(1)
	}

	for _, filePath := range os.Args[1:] {
		if err := sendFile(cfg, filePath); err != nil {
			fmt.Printf("\033[31m❌\033[0m %s: %v\n", filePath, err)
		}
	}
}
