package validator

type NameValidator interface {
	ValidateName(name string) error
}

type EmailValidator interface {
	ValidateEmail(email string) error
}

type PasswordValidator interface {
	ValidatePassword(password string) error
}
