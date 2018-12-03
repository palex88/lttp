package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	uuid2 "github.com/google/uuid"
)

func main() {

	configuration := parseConfigFile("config.json")
	fmt.Println(configuration.username)

	var (
		id        string
		email     string
		firstname string
		lastname  string
	)

	//sql := fmt.Sprintf("%s:%stcp(%s:3306)/%s", configuration.username, configuration.password, configuration.endpoint, configuration.database)

	db, err := sql.Open(
		"mysql",
		"")
	if err != nil {
		panic(err)
	}

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &email, &firstname, &lastname)
		if err != nil {
			panic(err)
		}
		fmt.Println(id, email, firstname, lastname)
	}

	defer db.Close()
}

func createUserId() {
	uuid, err := uuid2.NewUUID()
	if err != nil {
		panic(err)
	}

	uuidString := uuid.String()
	fmt.Println(uuidString)
	fmt.Println(len(uuidString))
}
