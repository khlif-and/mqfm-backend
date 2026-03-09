package helper

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"mqfm-backend/internal/shared/logger"

	"go.uber.org/zap"
)

type AudioConverter struct{}

func NewAudioConverter() *AudioConverter {
	return &AudioConverter{}
}

func (c *AudioConverter) ConvertToOGG(inputPath string) (string, error) {
	ext := filepath.Ext(inputPath)
	outputPath := strings.TrimSuffix(inputPath, ext) + ".ogg"

	cmd := exec.Command("ffmpeg", "-i", inputPath, "-c:a", "libvorbis", "-q:a", "5", "-y", outputPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Error("ffmpeg convert failed", zap.String("output", string(output)), zap.Error(err))
		return "", fmt.Errorf("ffmpeg conversion failed: %w", err)
	}

	return outputPath, nil
}

func (c *AudioConverter) CreateClip(inputPath string, startSec, endSec int) (string, error) {
	duration := endSec - startSec
	ext := filepath.Ext(inputPath)
	dir := filepath.Dir(inputPath)
	clipFilename := fmt.Sprintf("clip_%s_%d_%d%s", GenerateUniqueFilename("clip"), startSec, endSec, ext)
	outputPath := filepath.Join(dir, clipFilename)

	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-ss", fmt.Sprintf("%d", startSec),
		"-t", fmt.Sprintf("%d", duration),
		"-c", "copy",
		"-y", outputPath,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Error("ffmpeg clip failed", zap.String("output", string(output)), zap.Error(err))
		return "", fmt.Errorf("ffmpeg clipping failed: %w", err)
	}

	return outputPath, nil
}
