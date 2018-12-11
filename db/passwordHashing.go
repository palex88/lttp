package db

import (
	"crypto/rand"
	"golang.org/x/crypto/scrypt"
	"io"
	"log"
)

const  (
	pwSaltBytes = 32
	pwHashBytes = 64
)

func hashAndSalt(password string) ([]byte, []byte) {
	salt := make([]byte, pwSaltBytes)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		log.Fatal(err)
	}

	hash, err := scrypt.Key([]byte(password), salt, 1<<14, 8, 1, pwHashBytes)
	if err != nil {
		log.Fatal(err)
	}

	return hash, salt
}