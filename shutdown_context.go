package nstd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

var (
	_ context.Context = (*ShutdownContext)(nil)

	// gracefulShutdownSignal defines the signals that will trigger a graceful shutdown.
	gracefulShutdownSignal = []os.Signal{
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	}
)

// ShutdownContext is a context that will be canceled when syscall.SIGINT or syscall.SIGTERM is received.
type ShutdownContext struct {
	context.Context
	cancel  context.CancelCauseFunc
	sigChan chan os.Signal
}

// NewShutdownContextWithCause creates a new context that will be canceled when syscall.SIGINT or syscall.SIGTERM is received
// Usually used with s.Wait() for a graceful shutdown
// ShutdownContext gives more detail on which signal makes the context Done
// Use NewShutdownContext for simpler method
func NewShutdownContextWithCause(ctx context.Context) (*ShutdownContext, context.CancelCauseFunc) {
	ctx, cancel := context.WithCancelCause(ctx)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, gracefulShutdownSignal...)

	return &ShutdownContext{
		Context: ctx,
		cancel:  cancel,
		sigChan: sigChan,
	}, cancel
}

// Wait waits for the context to be done or an error to be received from the error channel.
// It returns an error if the context is canceled due to a signal or if an error is received from the error channel.
// If the context is canceled due to a signal, it returns a shutdownCause error that indicates which signal triggered the shutdown.
func (s *ShutdownContext) Wait(errChan <-chan error) error {
	select {
	case sig := <-s.sigChan:
		signal.Stop(s.sigChan)
		close(s.sigChan)
		s.cancel(newShutdownCause(sig))
		return context.Cause(s)
	case <-s.Done():
		return context.Cause(s)
	case err := <-errChan:
		return err
	}
}

// NewShutdownContext creates a new context that will be canceled when syscall.SIGINT or syscall.SIGTERM is received
// Usually used with WaitUntil to wait for a graceful shutdown.
// Use NewShutdownContextWithCause for more detail on which signal triggers the graceful shutdown
func NewShutdownContext(ctx context.Context) (context.Context, context.CancelFunc) {
	return signal.NotifyContext(ctx, gracefulShutdownSignal...)
}

// Wait waits until the context is done or an error is received from the error channel.
// It returns nil if the context is canceled or if the error channel is closed without an error.
// If an error is received from the error channel, it returns that error.
func Wait(ctx context.Context, errCh <-chan error) error {
	if sc, ok := ctx.(*ShutdownContext); ok {
		return sc.Wait(errCh)
	}

	select {
	case <-ctx.Done():
		return context.Cause(ctx)
	case err := <-errCh:
		return err
	}
}
