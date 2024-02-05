package security

import "golang.org/x/crypto/bcrypt"

type HashGenerator struct{}

func NewHashGenerator() *HashGenerator {
	return &HashGenerator{}
}

func (g HashGenerator) GenerateFromPassword(password string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(result), nil
}
