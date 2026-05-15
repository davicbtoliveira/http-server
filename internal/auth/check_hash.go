package auth

import "github.com/alexedwards/argon2id"

func CheckPasswordHash(password, hash string) (bool, error) {
	isHashed, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}

	return isHashed, err
}
