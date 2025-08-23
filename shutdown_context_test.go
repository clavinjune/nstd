package nstd_test

import (
	"context"
	"database/sql"
	"os"
	"syscall"
	"testing"
	"time"

	. "github.com/clavinjune/nstd"
)

func TestShutdownContext(t *testing.T) {
	t.Run("shutdown due to deadline", func(t *testing.T) {
		baseCtx, baseCancel := context.WithTimeout(t.Context(), time.Millisecond)
		defer baseCancel()

		ctx, cancel := NewShutdownContext(baseCtx)
		defer cancel()

		err := Wait(ctx, nil)
		RequireNotNil(t, err)
		RequireEqual(t, context.DeadlineExceeded, err)

	})

	t.Run("shutdown due to context cancelled", func(t *testing.T) {
		ctx, cancel := NewShutdownContext(t.Context())
		defer cancel()

		go func() {
			cancel()
		}()

		err := Wait(ctx, nil)
		RequireNotNil(t, err)
		RequireEqual(t, context.Canceled, err)
	})

	t.Run("shutdown due to syscall.SIGINT", func(t *testing.T) {
		ctx, cancel := NewShutdownContext(t.Context())
		defer cancel()

		go func() {
			RequireNil(t, syscall.Kill(os.Getpid(), syscall.SIGINT))
		}()

		err := Wait(ctx, nil)
		RequireNotNil(t, err)
		RequireEqual(t, context.Canceled, err)
	})

	t.Run("shutdown due to syscall.SIGTERM", func(t *testing.T) {
		ctx, cancel := NewShutdownContext(t.Context())
		defer cancel()

		go func() {
			RequireNil(t, syscall.Kill(os.Getpid(), syscall.SIGTERM))
		}()

		err := Wait(ctx, nil)
		RequireNotNil(t, err)
		RequireEqual(t, context.Canceled, err)
	})

	t.Run("ignore syscall.SIGCHLD", func(t *testing.T) {
		baseCtx, baseCancel := context.WithTimeout(t.Context(), time.Millisecond)
		defer baseCancel()

		ctx, cancel := NewShutdownContext(baseCtx)
		defer cancel()

		RequireNil(t, syscall.Kill(os.Getpid(), syscall.SIGCHLD))

		err := Wait(ctx, nil)
		RequireNotNil(t, err)
		RequireEqual(t, context.DeadlineExceeded, err)
	})

	t.Run("shutdown due to errorChan is filled", func(t *testing.T) {
		ctx, cancel := NewShutdownContext(t.Context())
		defer cancel()
		errChan := make(chan error, 1)

		go func() {
			errChan <- sql.ErrNoRows
		}()

		err := Wait(ctx, errChan)
		RequireNotNil(t, err)
		RequireEqual(t, sql.ErrNoRows, err)
	})
}

func TestShutdownContextWithCause(t *testing.T) {
	t.Run("shutdown due to deadline", func(t *testing.T) {
		baseCtx, baseCancel := context.WithTimeout(t.Context(), time.Millisecond)
		defer baseCancel()

		ctx, cancel := NewShutdownContextWithCause(baseCtx)
		defer cancel(context.Canceled)

		err := ctx.Wait(nil)
		RequireNotNil(t, err)
		RequireEqual(t, context.DeadlineExceeded, err)

	})

	t.Run("shutdown due to context cancelled", func(t *testing.T) {
		ctx, cancel := NewShutdownContextWithCause(t.Context())
		defer cancel(context.Canceled)

		go func() {
			cancel(context.Canceled)
		}()

		err := ctx.Wait(nil)
		RequireNotNil(t, err)
		RequireEqual(t, context.Canceled, err)
	})

	t.Run("shutdown due to syscall.SIGINT", func(t *testing.T) {
		ctx, cancel := NewShutdownContextWithCause(t.Context())
		defer cancel(context.Canceled)

		go func() {
			RequireNil(t, syscall.Kill(os.Getpid(), syscall.SIGINT))
		}()

		err := ctx.Wait(nil)
		RequireNotNil(t, err)
		RequireErrIs(t, err, context.Canceled)
	})

	t.Run("shutdown due to syscall.SIGTERM", func(t *testing.T) {
		ctx, cancel := NewShutdownContextWithCause(t.Context())
		defer cancel(context.Canceled)

		go func() {
			RequireNil(t, syscall.Kill(os.Getpid(), syscall.SIGTERM))
		}()

		err := ctx.Wait(nil)
		RequireNotNil(t, err)
		RequireErrIs(t, err, context.Canceled)
	})

	t.Run("ignore syscall.SIGCHLD", func(t *testing.T) {
		baseCtx, baseCancel := context.WithTimeout(t.Context(), time.Millisecond)
		defer baseCancel()

		ctx, cancel := NewShutdownContextWithCause(baseCtx)
		defer cancel(context.Canceled)

		RequireNil(t, syscall.Kill(os.Getpid(), syscall.SIGCHLD))

		err := ctx.Wait(nil)
		RequireNotNil(t, err)
		RequireEqual(t, context.DeadlineExceeded, err)
	})

	t.Run("shutdown due to errorChan is filled", func(t *testing.T) {
		ctx, cancel := NewShutdownContextWithCause(t.Context())
		defer cancel(context.Canceled)
		errChan := make(chan error, 1)

		go func() {
			errChan <- sql.ErrNoRows
		}()

		err := ctx.Wait(errChan)
		RequireNotNil(t, err)
		RequireEqual(t, sql.ErrNoRows, err)
	})
}

func TestShutdownContextCombined(t *testing.T) {
	t.Run("shutdown due to deadline", func(t *testing.T) {
		baseCtx, baseCancel := context.WithTimeout(t.Context(), time.Millisecond)
		defer baseCancel()

		ctx, cancel := NewShutdownContextWithCause(baseCtx)
		defer cancel(context.Canceled)

		err := Wait(ctx, nil)
		RequireNotNil(t, err)
		RequireEqual(t, context.DeadlineExceeded, err)

	})

	t.Run("shutdown due to context cancelled", func(t *testing.T) {
		ctx, cancel := NewShutdownContextWithCause(t.Context())
		defer cancel(context.Canceled)

		go func() {
			cancel(context.Canceled)
		}()

		err := Wait(ctx, nil)
		RequireNotNil(t, err)
		RequireEqual(t, context.Canceled, err)
	})

	t.Run("shutdown due to syscall.SIGINT", func(t *testing.T) {
		ctx, cancel := NewShutdownContextWithCause(t.Context())
		defer cancel(context.Canceled)

		go func() {
			RequireNil(t, syscall.Kill(os.Getpid(), syscall.SIGINT))
		}()

		err := Wait(ctx, nil)
		RequireNotNil(t, err)
		RequireErrIs(t, err, context.Canceled)
	})

	t.Run("shutdown due to syscall.SIGTERM", func(t *testing.T) {
		ctx, cancel := NewShutdownContextWithCause(t.Context())
		defer cancel(context.Canceled)

		go func() {
			RequireNil(t, syscall.Kill(os.Getpid(), syscall.SIGTERM))
		}()

		err := Wait(ctx, nil)
		RequireNotNil(t, err)
		RequireErrIs(t, err, context.Canceled)
	})

	t.Run("ignore syscall.SIGCHLD", func(t *testing.T) {
		baseCtx, baseCancel := context.WithTimeout(t.Context(), time.Millisecond)
		defer baseCancel()

		ctx, cancel := NewShutdownContextWithCause(baseCtx)
		defer cancel(context.Canceled)

		RequireNil(t, syscall.Kill(os.Getpid(), syscall.SIGCHLD))

		err := Wait(ctx, nil)
		RequireNotNil(t, err)
		RequireEqual(t, context.DeadlineExceeded, err)
	})

	t.Run("shutdown due to errorChan is filled", func(t *testing.T) {
		ctx, cancel := NewShutdownContextWithCause(t.Context())
		defer cancel(context.Canceled)
		errChan := make(chan error, 1)

		go func() {
			errChan <- sql.ErrNoRows
		}()

		err := Wait(ctx, errChan)
		RequireNotNil(t, err)
		RequireEqual(t, sql.ErrNoRows, err)
	})
}
