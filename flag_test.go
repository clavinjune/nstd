package nstd_test

import (
	"flag"
	"os"
	"testing"

	. "github.com/clavinjune/nstd"
)

func TestFlagSet_FromEnv(t *testing.T) {
	defer os.Clearenv()
	t.Setenv("TEST_STR", "from-env")
	t.Setenv("TEST_INT", "42")
	t.Setenv("TEST_BOOL", "true")

	fs := NewFlagSet("test", flag.ExitOnError)
	strFlag := fs.String("str", "default", "usage")
	intFlag := fs.Int("int", 0, "usage")
	boolFlag := fs.Bool("bool", false, "usage")

	RequireNil(t, fs.Parse("--str", "from-args", "--int", "100", "--bool"))
	RequireEqual(t, *strFlag, "from-env")
	RequireEqual(t, *intFlag, 42)
	RequireEqual(t, *boolFlag, true)
}

func TestFlagSet_FromArgs(t *testing.T) {
	defer os.Clearenv()

	fs := NewFlagSet("test", flag.ExitOnError)
	strFlag := fs.String("str", "default", "usage")
	intFlag := fs.Int("int", 0, "usage")
	boolFlag := fs.Bool("bool", false, "usage")

	RequireNil(t, fs.Parse("--str", "from-args", "--int", "100", "--bool"))
	RequireEqual(t, *strFlag, "from-args")
	RequireEqual(t, *intFlag, 100)
	RequireEqual(t, *boolFlag, true)
}

func TestFlagSet(t *testing.T) {
	tt := []struct {
		_        struct{}
		Name     string
		EnvKey   string
		EnvValue string
		Args     []string
		Want     string
	}{
		{
			Name:     "default value, not setting env or args",
			EnvKey:   "",
			EnvValue: "",
			Args:     nil,
			Want:     "default",
		},
		{
			Name:     "setting correct env",
			EnvKey:   "TEST_NAME",
			EnvValue: "from-env",
			Args:     nil,
			Want:     "from-env",
		},
		{
			Name:     "setting not related env",
			EnvKey:   "TEST_NAME_NOT_RELATED",
			EnvValue: "not-related",
			Args:     nil,
			Want:     "default",
		},
		{
			Name:     "setting correct args",
			EnvKey:   "",
			EnvValue: "",
			Args:     []string{"--name", "from-args"},
			Want:     "from-args",
		},
		{
			Name:     "setting both correct args and env",
			EnvKey:   "TEST_NAME",
			EnvValue: "from-env",
			Args:     []string{"--name", "from-args"},
			Want:     "from-env",
		},
	}

	for i := range tt {
		tc := tt[i]

		t.Run(tc.Name, func(t *testing.T) {
			defer os.Clearenv()
			if tc.EnvKey != "" && tc.EnvValue != "" {
				t.Setenv(tc.EnvKey, tc.EnvValue)
			}

			fs := NewFlagSet("test", flag.ExitOnError)
			nameFlag := fs.String("name", "default", "usage")

			RequireNotNil(t, fs.FlagSet())
			RequireNil(t, fs.Parse(tc.Args...))
			RequireEqual(t, *nameFlag, tc.Want)
		})
	}
}
