package hash

import (
	"crypto/sha1"
	"fmt"
)

type PasswordHashes interface {
	Hash(password string) (string, error)
}

type SHA1Hashes struct {
	salt string
}

func NewSHA1Hashes(salt string) *SHA1Hashes {
	return &SHA1Hashes{salt: salt}
}

func (h *SHA1Hashes) Hash(password string) (string, error) {
	hash := sha1.New()

	if _, err := hash.Write([]byte(password)); err != nil {
		return "", fmt.Errorf("hash.Hash: %w", err)
	}

	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}
