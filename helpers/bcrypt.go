package helpers

import "golang.org/x/crypto/bcrypt"

type IHasher interface {
	HashString(plainText string) (string, error)
	CheckHash(hashedText string, plainText string) bool
}

type Hasher struct{}

func NewHasher() IHasher {
	return &Hasher{}
}

func (h *Hasher) HashString(plainText string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(plainText), 14)
	return string(hashedBytes), err
}

func (h *Hasher) CheckHash(hashedText string, plainText string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedText), []byte(plainText))
	return err == nil
}
