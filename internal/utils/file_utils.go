package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateUniqueFilename generates a unique filename to prevent overwrites.
func GenerateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	filename := strings.TrimSuffix(originalFilename, ext)
	
	// Create a safe filename (remove spaces, etc)
	filename = strings.ReplaceAll(filename, " ", "_")
	
	// Add UUID and Timestamp
	uniqueID := uuid.New().String()
	timestamp := time.Now().Unix()
	
	return filename + "_" + uniqueID + "_" + string(rune(timestamp)) + ext
}

// SaveUploadedFile saves a multipart file to the destination path.
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return err
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

// DownloadImage downloads an image from a URL, saves it to destDir with a unique filename, and returns the relative path.
func DownloadImage(url, destDir string) (string, error) {
	if url == "" {
		return "", nil
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image: status %d", resp.StatusCode)
	}

	ext := ".jpg"
	if strings.Contains(resp.Header.Get("Content-Type"), "png") {
		ext = ".png"
	}

	filename := GenerateUniqueFilename("profile" + ext)
	dst := filepath.Join(destDir, filename)

	if err = os.MkdirAll(filepath.Dir(dst), 0750); err != nil {
		return "", err
	}

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	// Always store and return forward-slash paths for web compatibility
	return filepath.ToSlash(dst), nil
}
