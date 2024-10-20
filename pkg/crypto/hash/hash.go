package hash

import (
	"crypto/sha1"
	"fmt"
)

type Hasher interface {
	Hash(s string) string
}

type sha1Hasher struct {
	Salt []byte
}

func NewSHA1Hasher(salt []byte) Hasher {
	return &sha1Hasher{Salt: salt}
}

func (h *sha1Hasher) Hash(s string) string {
	hash := sha1.New()
	hash.Write([]byte(s))
	return fmt.Sprintf("%x", hash.Sum(h.Salt))
}
