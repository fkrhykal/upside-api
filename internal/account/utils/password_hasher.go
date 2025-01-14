package utils

type PasswordHasher interface {
	Hash(rawPassword string) (string, error)
	Match(rawPassword string, hashedPassword string) bool
}
