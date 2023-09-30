package logtest

import (
	"bytes"
	"log/slog"
	"testing"
)

func NewLogger(t *testing.T, level slog.Level) (*slog.Logger, *bytes.Buffer) {
	var buf bytes.Buffer
	opts := &slog.HandlerOptions{
		Level: level,
	}
	h := slog.NewTextHandler(&buf, opts)
	ll := slog.New(h)

	return ll, &buf
}

func NewDefaultLogger(t *testing.T) (*slog.Logger, *bytes.Buffer) {
	return NewLogger(t, slog.LevelDebug)
}
