package errb

import "fmt"

type Builder interface {
	Errorf(format string, args ...any) error
}

func Errorf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}
