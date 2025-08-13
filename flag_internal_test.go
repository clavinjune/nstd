package nstd

import (
	"flag"
	"os"
	"testing"
)

func TestFlagSet_getFromEnv(t *testing.T) {
	defer os.Clearenv()
	t.Setenv("TEST_STR", "from-env")
	t.Setenv("TEST_INT", "42")
	t.Setenv("TEST_BOOL", "true")

	fs := NewFlagSet("test", flag.ExitOnError)

	val, ok := fs.getFromEnv("str")
	RequireTrue(t, ok)
	RequireEqual(t, "from-env", val)
	val, ok = fs.getFromEnv("int")
	RequireTrue(t, ok)
	RequireEqual(t, "42", val)
	val, ok = fs.getFromEnv("bool")
	RequireTrue(t, ok)
	RequireEqual(t, "true", val)
}
