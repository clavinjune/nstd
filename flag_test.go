package nstd_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/clavinjune/nstd"
)

func ExampleNewFlagSet() {
	defer os.Clearenv()
	os.Setenv("EXAMPLE_NAME", "from-envs")
	fs := nstd.NewFlagSet("example", flag.ExitOnError)
	nameFlag := fs.String("name", "default", "usage")

	if err := fs.Parse([]string{"-name=from-args"}...); err != nil {
		panic(err)
	}

	fmt.Println(*nameFlag)
	// Output: from-envs
}

func TestFlagSet_String(t *testing.T) {
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

			fs := nstd.NewFlagSet("test", flag.ExitOnError)
			nameFlag := fs.String("name", "default", "usage")

			if err := fs.Parse(tc.Args...); err != nil {
				t.Fatalf("error on fs.Parse: %+q", err.Error())
			}

			if *nameFlag != tc.Expected {
				t.Fatalf("expected: %+q, actual: %+q", tc.Expected, *nameFlag)
			}

		})
	}
}
