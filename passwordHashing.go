package main

import (
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
	"log"
)

const  (
	pwHashBytes = 64
)

func hashAndSalt(password string) []byte {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Fatal(err)
	}

	return hash
}

func checkPassword(password []byte, salt []byte) ([]byte, error) {
	hash, err := scrypt.Key(password, salt, 1<<14, 8, 1, pwHashBytes)
	return hash, err
}
