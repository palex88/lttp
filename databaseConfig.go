package main

import (
	"database/sql"
	"fmt"
	"log"
)

var Database *sql.DB

func configureDatabase(config Config) (db *sql.DB, err error) {

	conn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
		config.Username,
		config.Password,
		config.Endpoint,
		config.Database)

	db, err = sql.Open("mysql", conn)
	if err != nil {
		log.Println(err)
	}

	return db, err
}
