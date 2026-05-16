// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package tests

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/llingr/llingr-nexus/nexus"
)

func TestDefaultLogger(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug, // capture all levels
	})

	logger := &nexus.DefaultLogger{
		SLog: slog.New(handler),
	}

	ctx := context.Background()

	t.Run("Error logs message", func(t *testing.T) {
		buf.Reset()
		logger.Error(ctx, "error message")
		if !strings.Contains(buf.String(), "error message") {
			t.Errorf("expected 'error message' in output, got: %s", buf.String())
		}
	})

	t.Run("Warn logs message", func(t *testing.T) {
		buf.Reset()
		logger.Warn(ctx, "warn message")
		if !strings.Contains(buf.String(), "warn message") {
			t.Errorf("expected 'warn message' in output, got: %s", buf.String())
		}
	})

	t.Run("Info logs message", func(t *testing.T) {
		buf.Reset()
		logger.Info(ctx, "info message")
		if !strings.Contains(buf.String(), "info message") {
			t.Errorf("expected 'info message' in output, got: %s", buf.String())
		}
	})

	t.Run("Debug logs message", func(t *testing.T) {
		buf.Reset()
		logger.Debug(ctx, "debug message")
		if !strings.Contains(buf.String(), "debug message") {
			t.Errorf("expected 'debug message' in output, got: %s", buf.String())
		}
	})
}

func TestNewDefaultLogger(t *testing.T) {
	logger := nexus.NewDefaultLogger(slog.LevelInfo)
	if logger == nil {
		t.Error("expected non-nil logger")
	}
}
