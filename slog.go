package nstd

import (
	"io"
	"log/slog"
)

// NewSlog simplifies the creation of a slog.Logger with options for debug level and structured logging.
func NewSlog(w io.Writer, isDebug, isStructured bool) *slog.Logger {
	l := slog.LevelInfo
	if isDebug {
		l = slog.LevelDebug
	}
	ho := &slog.HandlerOptions{
		AddSource: isDebug,
		Level:     l,
	}

	var h slog.Handler
	if isStructured {
		h = slog.NewJSONHandler(w, ho)
	} else {
		h = slog.NewTextHandler(w, ho)
	}

	return slog.New(h)
}
