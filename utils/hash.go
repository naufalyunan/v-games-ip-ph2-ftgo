package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password []byte) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		return nil, err
	}
	return hashed, nil
}
