package nstd

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestShutdownCause(t *testing.T) {
	err := newShutdownCause(os.Interrupt)
	fmt.Println(err.Error())
	fmt.Println(err.String())
	RequireNotNil(t, err)
	RequireNotNil(t, err)
	RequireErrIs(t, err, context.Canceled)
	var sc *shutdownCause
	RequireErrAs(t, err, &sc)
	sc.Signal() // Ensure it implements os.Signal interface
}
