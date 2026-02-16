package cryptx

import (
	"golang.org/x/crypto/bcrypt"
)

func PasswordEncrypt(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(dbHash, inputPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbHash), []byte(inputPassword))
	return err == nil
}
