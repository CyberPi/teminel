package flag

import (
	"strings"
)

type MultiFlag []string

func (flag *MultiFlag) String() string {
	return strings.Join(*flag, ",")
}

func (flag *MultiFlag) Set(value string) error {
	*flag = append(*flag, value)
	return nil
}

func (flag *MultiFlag) Default(defaults ...string) error {
	if len(*flag) == 0 {
		*flag = defaults
	}
	return nil
}
