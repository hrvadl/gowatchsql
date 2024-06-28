package logger

import (
	"io"
	"log/slog"
)

func New(to io.Writer) *slog.Logger {
	return slog.New(slog.NewTextHandler(to, nil))
}
