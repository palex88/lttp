package main

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	uuid2 "github.com/google/uuid"
)

var Conn *sql.DB

func init() {

	var err error

	config := ParseConfigs()

	conn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
		config.Username,
		config.Password,
		config.Endpoint,
		config.Database)

	Conn, err = sql.Open("mysql", conn)
	if err != nil {
		log.Println(err)
	}
}

func CreateUserId() (uuidStr string) {
	uuid, err := uuid2.NewUUID()
	if err != nil {
		panic(err)
	}

	return uuid.String()
}

func CreateUser(email string, firstName string, lastName string, password string) (result sql.Result, err error) {

	userId := CreateUserId()
	hashedPassword := hashAndSalt(password)

	query := fmt.Sprintf(
		"INSERT INTO users (id, email, firstname, lastname, hashedpassword) VALUES ('%s', '%s', '%s', '%s', '%s')",
		userId, email, firstName, lastName, string(hashedPassword))

	result, err = Conn.Exec(query)
	if err != nil {
		fmt.Println(err)
	}

	return result, err
}

func AuthUser(email string, password string) bool {

	var hashedPassword string
	query := fmt.Sprintf("SELECT hashedpassword FROM users WHERE email='%s' LIMIT 1", email)

	rows, err := Conn.Query(query)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&hashedPassword)
		if err != nil {
			log.Println(err)
		}
	}

	bPassword := []byte(password)
	bHash := []byte(hashedPassword)

	err = bcrypt.CompareHashAndPassword(bHash, bPassword)
	log.Println("PW Compare: ", err)
	if err == nil {
		return true
	} else {
		return false
	}
}

func GetAllUsers() (allRows []User, err error) {

	var (
		id        string
		email     string
		firstname string
		lastname  string
	)

	rows, err := Conn.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &email, &firstname, &lastname)
		if err != nil {
			fmt.Println(err)
		}
		u := User{id, email, firstname, lastname}
		allRows = append(allRows, u)
	}

	return allRows, err
}

func GetAllLinks(userId string) (allLinks []string, err error) {

	var link string

	query := fmt.Sprintf("SELECT link FROM links WHERE userid='%s'", userId)

	rows, err := Conn.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&link)
		if err != nil {
			fmt.Print(err)
		}

		allLinks = append(allLinks, link)
	}

	return allLinks, err
}

func GetUserId(email string) (id string, err error) {

	query := fmt.Sprintf("SELECT id FROM users WHERE email='%s'", email)

	rows, err := Conn.Query(query)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			fmt.Println(err)
		}
	}

	return id, err
}
