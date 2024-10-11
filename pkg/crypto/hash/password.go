package hash

import (
	"crypto/sha1"
	"fmt"
)

type PasswordHasher interface {
	Hash(password string) string
}

type PasswordHasherSHA1 struct {
	salt []byte
}

func NewPasswordHasher(salt []byte) *PasswordHasherSHA1 {
	return &PasswordHasherSHA1{salt: salt}
}

func (h *PasswordHasherSHA1) Hash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(h.salt))
}
