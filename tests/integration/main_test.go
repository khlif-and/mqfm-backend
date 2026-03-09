package integration_test

import (
	"os"
	"testing"

	"mqfm-backend/internal/shared/logger"
)

func TestMain(m *testing.M) {
	logger.Init()
	os.Exit(m.Run())
}
