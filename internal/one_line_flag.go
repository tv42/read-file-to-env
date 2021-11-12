package internal

import (
	"errors"
	"flag"
	"strings"
)

type EnvVarFromFile struct {
	Name string
	Path string
}

type OneLineFlag []EnvVarFromFile

var _ flag.Value = (*OneLineFlag)(nil)

func (o *OneLineFlag) String() string {
	// Multiple flags don't really coalesce back to one string!
	return "..."
}

func (o *OneLineFlag) Set(input string) error {
	idx := strings.IndexByte(input, '=')
	if idx < 0 {
		return errors.New("value must have an equals sign")
	}
	name := input[:idx]
	if name == "" {
		return errors.New("environment variable name missing")
	}
	if strings.IndexByte(name, '\x00') >= 0 {
		return errors.New("environment variable name cannot contain NUL")
	}
	p := input[idx+1:]
	(*o) = append(*o, EnvVarFromFile{Name: name, Path: p})
	return nil
}
