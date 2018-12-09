package main

import (
	"database/sql"
	"fmt"
)

func main() {

	var db *sql.DB
	config := ParseConfigs("config.json")

	db = OpenDatabase(config)

	allRows := GetAllUsers(db)
	for _, element := range allRows {
		fmt.Println(element.Email)
	}

	allLinks := GetAllLinks(db, "60049ce6-f6c3-11e8-aaa6-dca9047d1371")
	for _, element := range allLinks {
		fmt.Println(element)
	}

	//result := CreateUser(db, "alex@test.com", "Alex", "Faker")
	//fmt.Println(result)

	//var port = ":8080"
	//
	//http.HandleFunc("/view/", makeHandler(viewHandler))
	//http.HandleFunc("/edit/", makeHandler(editHandler))
	//http.HandleFunc("/save/", makeHandler(saveHandler))
	//
	//log.Fatal(http.ListenAndServe(port, nil))
}
