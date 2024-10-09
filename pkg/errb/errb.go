package errb

import "fmt"

type Builder interface {
	Errorf(format string, args ...any) error
}

type stdBuilderImpl struct{}

func (b *stdBuilderImpl) Errorf(format string, args ...any) error {
	return fmt.Errorf(format, args...)
}

var stdBuilder Builder = &stdBuilderImpl{}

func Default() Builder {
	return stdBuilder
}

func SetDefault(b Builder) {
	stdBuilder = b
}

func Errorf(format string, args ...any) error {
	return stdBuilder.Errorf(format, args...)
}
