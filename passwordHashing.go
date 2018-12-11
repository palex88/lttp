package main

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

func checkPassword(password []byte, salt []byte) ([]byte, error) {
	hash, err := scrypt.Key(password, salt, 1<<14, 8, 1, pwHashBytes)
	return hash, err
}
