package application

type PasswordHasher interface {
	Hash(string) (string, error)
}
