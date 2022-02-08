package helper

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	passwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(passwd), err
}

func CheckPasswordHash(hashVal, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashVal), []byte(password))
	if err != nil {
		return false
	} else {
		return true
	}
}
