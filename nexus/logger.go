// SPDX-FileCopyrightText: Copyright (c) 2025 The llingr-nexus Authors
// SPDX-License-Identifier: Apache-2.0

package nexus

import (
	"context"
	"fmt"
	"log/slog"
	"os"
)

// Logger simple facade with zero dependencies;
// replace with any back-end.
type Logger interface {
	Error(ctx context.Context, format string, args ...any)

	Warn(ctx context.Context, format string, args ...any)

	Info(ctx context.Context, format string, args ...any)

	Debug(ctx context.Context, format string, args ...any)
}

// NewDefaultLogger implementation using SLog
func NewDefaultLogger(level slog.Level) Logger {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
	})
	return &DefaultLogger{
		SLog: slog.New(handler),
		ctx:  context.Background(), // pre-allocate fixed context
	}
}

// DefaultLogger for basic functionality - production deployments
// can use different logging back-ends, for example uber-go/zap.
//
// This implementation uses printf semantics.
type DefaultLogger struct {
	SLog *slog.Logger
	ctx  context.Context // to avoid unnecessary GC pressure
}

// Error logs at error level. Context parameter satisfies Logger interface;
// this default implementation uses a pre-allocated context to avoid per-call allocation.
// Production deployments should use a Logger implementation with proper context propagation.
//
//nolint:contextcheck // deliberate: default logger uses pre-allocated context for zero-alloc logging
func (l *DefaultLogger) Error(_ context.Context, format string, args ...any) {
	l.SLog.ErrorContext(l.ctx, fmt.Sprintf(format, args...))
}

// Warn logs at warn level. See Error for context handling rationale.
//
//nolint:contextcheck // deliberate: default logger uses pre-allocated context for zero-alloc logging
func (l *DefaultLogger) Warn(_ context.Context, format string, args ...any) {
	l.SLog.WarnContext(l.ctx, fmt.Sprintf(format, args...))
}

// Info logs at info level. See Error for context handling rationale.
//
//nolint:contextcheck // deliberate: default logger uses pre-allocated context for zero-alloc logging
func (l *DefaultLogger) Info(_ context.Context, format string, args ...any) {
	l.SLog.InfoContext(l.ctx, fmt.Sprintf(format, args...))
}

// Debug logs at debug level. See Error for context handling rationale.
//
//nolint:contextcheck // deliberate: default logger uses pre-allocated context for zero-alloc logging
func (l *DefaultLogger) Debug(_ context.Context, format string, args ...any) {
	l.SLog.DebugContext(l.ctx, fmt.Sprintf(format, args...))
}
