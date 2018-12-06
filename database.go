package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	uuid2 "github.com/google/uuid"
)

func CreateUserId() {
	uuid, err := uuid2.NewUUID()
	if err != nil {
		panic(err)
	}

	uuidString := uuid.String()
	fmt.Println(uuidString)
	fmt.Println(len(uuidString))
}

func OpenDatabaseConnection(config Config) {

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

	fmt.Printf("Conn: %s\n", conn)

	db, err := sql.Open("mysql", conn)
	fmt.Println("DB: ", db)
	if err != nil {
		fmt.Println("OPEN FAILURE")
		fmt.Println(err)
	}
	defer db.Close()

	fmt.Println("getting rows")
	rows, err := db.Query("SELECT * FROM users")
	fmt.Println("rows: ", rows)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("rows")
	fmt.Println(rows)
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &email, &firstname, &lastname)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(id, email, firstname, lastname)
	}
}

func main() {
	config := ParseConfigs("config.json")
	fmt.Printf("Username: %s\n", config.Username)
	fmt.Printf("Password: %s\n", config.Password)
	fmt.Printf("Endpoint: %s\n", config.Endpoint)
	fmt.Printf("Database: %s\n", config.Database)

	OpenDatabaseConnection(config)
}
