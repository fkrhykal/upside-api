package utils

import "golang.org/x/crypto/bcrypt"

type BcryptPasswordHasher struct {
}

func (b *BcryptPasswordHasher) Hash(rawPassword string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func (b *BcryptPasswordHasher) Match(rawPassword string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(rawPassword))
	return err == nil
}

func NewBcryptPasswordHasher() PasswordHasher {
	return &BcryptPasswordHasher{}
}
