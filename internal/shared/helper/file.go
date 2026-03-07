package helper

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func GenerateUniqueFilename(originalFilename string) string {
	ext := filepath.Ext(originalFilename)
	name := strings.TrimSuffix(originalFilename, ext)
	name = strings.ReplaceAll(name, " ", "_")
	return fmt.Sprintf("%s_%s%s", name, uuid.New().String(), ext)
}

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
		return "", fmt.Errorf("download failed: status %d", resp.StatusCode)
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

	return filepath.ToSlash(dst), nil
}
