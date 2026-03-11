package audiomock

import "mqfm-backend/internal/domain/port"

type MockColorExtractorService struct {
	ExtractFn func(imagePath string) (string, error)
}

func (m *MockColorExtractorService) ExtractDominantColor(imagePath string) (string, error) {
	return m.ExtractFn(imagePath)
}

var _ port.ColorExtractorService = (*MockColorExtractorService)(nil)
