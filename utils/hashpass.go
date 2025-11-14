package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (encryptedPass string, err error) {

	salt := make([]byte, 16)

	_, err = rand.Read(salt)

	if err != nil {
		return "", err
	}

	hashPass := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	saltBase64 := base64.StdEncoding.EncodeToString(salt)
	hashPassBase64 := base64.StdEncoding.EncodeToString(hashPass)

	encryptedPass = fmt.Sprintf("%s.%s", saltBase64, hashPassBase64)
	return encryptedPass, nil
}
