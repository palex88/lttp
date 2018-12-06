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

func CreateUserId() (uuidStr string) {
	uuid, err := uuid2.NewUUID()
	if err != nil {
		panic(err)
	}

	return uuid.String()
}

func OpenDatabaseConnection(config Config) (allRows []User) {

	var (
		id        string
		email     string
		firstname string
		lastname  string
	)

	conn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
		config.Username,
		config.Password,
		config.Endpoint,
		config.Database)

	db, err := sql.Open("mysql", conn)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

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

func main() {
	config := ParseConfigs("config.json")
	fmt.Printf("Username: %s\n", config.Username)
	fmt.Printf("Password: %s\n", config.Password)
	fmt.Printf("Endpoint: %s\n", config.Endpoint)
	fmt.Printf("Database: %s\n", config.Database)

	allRows := OpenDatabaseConnection(config)
	for _, element := range allRows {
		fmt.Println(element.Email)
	}
}
