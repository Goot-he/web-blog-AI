package utils

import (
	"encoding/base64"
	"log"

	"golang.org/x/crypto/scrypt"
)

func Bcrypt(password string) string {
	const keyLen = 10
	salt := []byte{12, 32, 4, 66, 66, 22, 222, 11}
	hashPassword, err := scrypt.Key([]byte(password), salt, 16384, 8, 1, keyLen)
	if err != nil {
		log.Fatal("bcrypt password failed, err:", err)
	}

	return base64.StdEncoding.EncodeToString(hashPassword)
}
