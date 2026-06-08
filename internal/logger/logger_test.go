package logger_test

import (
	"bytes"
	"garnet/internal/logger"
	"testing"
)

func TestLoggerInitialization(t *testing.T) {
	if logger.Logger == nil {
		t.Fatal("Expected logger to be initialized, got nil")
	}

	// Test that we can write to the logger
	var buf bytes.Buffer
	oldWriter := logger.Logger.Writer()
	defer logger.Logger.SetOutput(oldWriter)

	logger.Logger.SetOutput(&buf)
	logger.Logger.Print("test log")

	if buf.Len() == 0 {
		t.Error("Expected logger to write to buffer, but it was empty")
	}
}
