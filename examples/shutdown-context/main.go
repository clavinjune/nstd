package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/clavinjune/nstd"
)

func main() {
	ctx, cancel := nstd.NewShutdownContextWithCause(context.Background())
	defer cancel(context.Canceled)
	errChan := make(chan error, 1)

	go func() {
		slog.Info("starting server")
		if err := http.ListenAndServe(":0", nil); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			errChan <- err
		}
	}()

	if err := ctx.Wait(errChan); err != nil {
		slog.Error(err.Error())
	}
	// press Ctrl+C to trigger shutdown
}
