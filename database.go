package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	uuid2 "github.com/google/uuid"
)

type User struct {
	Id        string
	Email     string
	FirstName string
	LastName  string
}

func OpenDatabase(config Config) (db *sql.DB) {

	conn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
		config.Username,
		config.Password,
		config.Endpoint,
		config.Database)

	db, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println(err)
	}

	return db
}

func CreateUserId() (uuidStr string) {
	uuid, err := uuid2.NewUUID()
	if err != nil {
		panic(err)
	}

	return uuid.String()
}

func CreateUser(db *sql.DB, email string, firstName string, lastName string) (result sql.Result) {

	userId := CreateUserId()

	query := fmt.Sprintf(
		"INSERT INTO users (id, email, firstname, lastname) VALUES ('%s', '%s', '%s', '%s')",
		userId, email, firstName, lastName)

	result, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
	}

	return result
}

func GetAllUsers(db *sql.DB) (allRows []User) {

	var (
		id        string
		email     string
		firstname string
		lastname  string
	)

	rows, err := db.Query("SELECT * FROM users")
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

	return allRows
}

func GetAllLinks(db *sql.DB, userId string) (allLinks []string) {

	var link string

	query := fmt.Sprintf("SELECT link FROM links WHERE userid='%s'", userId)

	rows, err := db.Query(query)
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

	return allLinks
}
