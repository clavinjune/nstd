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

	RequireNoErr(t, fs.Parse("--str", "from-args", "--int", "100", "--bool"))
	RequireEqual(t, "from-env", *strFlag)
	RequireEqual(t, 42, *intFlag)
	RequireEqual(t, true, *boolFlag)
}

func TestFlagSet_FromArgs(t *testing.T) {
	defer os.Clearenv()

	fs := NewFlagSet("test", flag.ExitOnError)
	strFlag := fs.String("str", "default", "usage")
	intFlag := fs.Int("int", 0, "usage")
	boolFlag := fs.Bool("bool", false, "usage")

	RequireNoErr(t, fs.Parse("--str", "from-args", "--int", "100", "--bool"))
	RequireEqual(t, "from-args", *strFlag)
	RequireEqual(t, 100, *intFlag)
	RequireEqual(t, true, *boolFlag)
}

func TestFlagSet(t *testing.T) {
	tt := []struct {
		_        struct{}
		Name     string
		EnvKey   string
		EnvValue string
		Args     []string
		Expected string
	}{
		{
			Name:     "default value, not setting env or args",
			EnvKey:   "",
			EnvValue: "",
			Args:     nil,
			Expected: "default",
		},
		{
			Name:     "setting correct env",
			EnvKey:   "TEST_NAME",
			EnvValue: "from-env",
			Args:     nil,
			Expected: "from-env",
		},
		{
			Name:     "setting not related env",
			EnvKey:   "TEST_NAME_NOT_RELATED",
			EnvValue: "not-related",
			Args:     nil,
			Expected: "default",
		},
		{
			Name:     "setting correct args",
			EnvKey:   "",
			EnvValue: "",
			Args:     []string{"--name", "from-args"},
			Expected: "from-args",
		},
		{
			Name:     "setting both correct args and env",
			EnvKey:   "TEST_NAME",
			EnvValue: "from-env",
			Args:     []string{"--name", "from-args"},
			Expected: "from-env",
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
			RequireNoErr(t, fs.Parse(tc.Args...))
			RequireEqual(t, tc.Expected, *nameFlag)
		})
	}
}
