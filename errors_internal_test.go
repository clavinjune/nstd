package nstd

import (
	"context"
	"os"
	"testing"
)

func TestShutdownCause(t *testing.T) {
	err := newShutdownCause(os.Interrupt)
	RequireNotNil(t, err)
	RequireNotNil(t, err)
	RequireErrIs(t, err, context.Canceled)
	var sc *shutdownCause
	RequireErrAs(t, err, &sc)
	sc.Signal() // Ensure it implements os.Signal interface
}
