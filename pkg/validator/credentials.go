package validator

import (
	"regexp"
	"unicode"
)

type CredentialsValidator struct{}

func NewCredentialsValidator() *CredentialsValidator {
	return &CredentialsValidator{}
}

func (v *CredentialsValidator) ValidateName(name string) error {
	if err := checkInputLength(name, 3, 256); err != nil {
		return err
	}

	var containsUpper bool
	for _, r := range name {
		if unicode.IsLetter(r) {
			if unicode.IsUpper(r) {
				containsUpper = true
			}
		} else if !unicode.IsNumber(r) {
			if !unicode.IsSpace(r) {
				return errInvalidCharacters
			}
		}
	}

	if !containsUpper {
		return errMustContainUpper
	}

	return nil
}

func (v *CredentialsValidator) ValidateEmail(email string) error {
	if err := checkInputLength(email, 4, 256); err != nil {
		return err
	}

	regex, err := regexp.Compile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if err != nil {
		return err
	}

	if !regex.MatchString(email) {
		return errInvalidEmail(email)
	}

	return nil
}

func (v *CredentialsValidator) ValidatePassword(password string) error {
	if err := checkInputLength(password, 8, 256); err != nil {
		return err
	}

	var (
		containsUpper  bool
		containsNumber bool
	)
	for _, r := range password {
		if unicode.IsLetter(r) {
			if unicode.IsUpper(r) {
				containsUpper = true
			}
		} else if unicode.IsNumber(r) {
			containsNumber = true
		} else {
			if !unicode.IsSpace(r) {
				return errInvalidCharacters
			}
		}
	}

	if !containsUpper {
		return errMustContainUpper
	} else if !containsNumber {
		return errMustContainNumber
	}

	return nil
}

func checkInputLength(input string, min, max int) error {
	if l := len(input); l < min {
		return errInputLengthLesserThanMin(min)
	} else if l > max {
		return errInputLengthBiggerThanMax(max)
	}
	return nil
}
