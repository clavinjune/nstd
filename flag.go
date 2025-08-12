package nstd

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	// reSymbols is a regular expression that matches any non-word character.
	reSymbols *regexp.Regexp = regexp.MustCompile(`\W+`)
)

// FlagSet wraps flag.FlagSet to provide a structured way to manage command-line flags with environment variable support.
// environment variable is prioritized over command-line arguments.
type FlagSet struct {
	_    struct{}
	std  *flag.FlagSet
	name string
}

// NewFlagSet creates a new FlagSet with the given name and error handling mode.
// the name is used to construct environment variable names by appending the flag name.
// for example: if the name is "myapp", the environment variable for a flag named "port" would be "MYAPP_PORT".
func NewFlagSet(name string, errorHandling flag.ErrorHandling) *FlagSet {
	fs := flag.NewFlagSet(name, errorHandling)

	return &FlagSet{
		std:  fs,
		name: strings.TrimSpace(name),
	}
}

// FlagSet returns the underlying flag.FlagSet.
func (fs *FlagSet) FlagSet() *flag.FlagSet {
	return fs.std
}

// Parse wraps the standard flag.FlagSet Parse method.
func (fs *FlagSet) Parse(args ...string) error {
	return fs.std.Parse(args)
}

// Bool wraps the standard flag.FlagSet BoolVar method to support environment variables.
func (fs *FlagSet) Bool(name string, value bool, usage string) *bool {
	f := fs.std.Bool(name, value, usage)
	e, ok := fs.getFromEnv(name)
	if !ok {
		return f
	}

	if i, err := strconv.ParseBool(e); err != nil {
		panic(err.Error())
	} else {
		return &i
	}
}

// String wraps the standard flag.FlagSet StringVar method to support environment variables.
func (fs *FlagSet) String(name, value, usage string) *string {
	f := fs.std.String(name, value, usage)
	if e, ok := fs.getFromEnv(name); ok {
		return &e
	}

	return f
}

// Int wraps the standard flag.FlagSet IntVar method to support environment variables.
func (fs *FlagSet) Int(name string, value int, usage string) *int {
	f := fs.std.Int(name, value, usage)
	e, ok := fs.getFromEnv(name)
	if !ok {
		return f
	}

	if i, err := strconv.Atoi(e); err != nil {
		panic(err.Error())
	} else {
		return &i
	}
}

// getFromEnv retrieves the value of an environment variable constructed from the FlagSet name and the flag name.
// it replace all non-word characters in the flag name with underscores and converts it to uppercase.
func (fs *FlagSet) getFromEnv(name string) (string, bool) {
	return os.LookupEnv(
		reSymbols.ReplaceAllString(
			strings.ToUpper(
				fmt.Sprintf(
					"%s_%s",
					fs.name,
					strings.TrimSpace(name),
				),
			),
			"_",
		),
	)
}
