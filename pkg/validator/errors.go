package validator

import (
	"errors"
	"fmt"
)

var (
	errMustContainUpper  = errors.New("input must contain uppercase")
	errMustContainNumber = errors.New("input must contain numbers")
	errInvalidCharacters = errors.New("input must contain only letters and numbers")
)

func errInputLengthLesserThanMin(min int) error {
	return fmt.Errorf("input length must be bigger than %d", min)
}

func errInputLengthBiggerThanMax(max int) error {
	return fmt.Errorf("input length must be lesser than %d", max)
}

func errInvalidEmail(email string) error {
	return fmt.Errorf("invalid email %s", email)
}
