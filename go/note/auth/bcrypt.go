package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func GenBcrypt(pwd []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(pwd, bcrypt.DefaultCost)
}

func ValidateBcrypt(pwd, hashed []byte) (ok bool, err error) {
	if err = bcrypt.CompareHashAndPassword(hashed, pwd); err != nil {
		return false, err
	}

	return true, nil
}
