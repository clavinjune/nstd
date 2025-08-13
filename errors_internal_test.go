package nstd

import (
	"context"
	"os"
	"testing"
)

func TestShutdownCause(t *testing.T) {
	err := newShutdownCause(os.Interrupt)
	RequireErrIs(t, err, context.Canceled)
	var sc *shutdownCause
	RequireErrAs(t, err, &sc)

}
