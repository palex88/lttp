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

func AuthUser(email string, password string) (user User, ok bool) {

	var (
		id string
		firstname string
		lastname string
		hashedPassword string
	)

	query := fmt.Sprintf("SELECT id, firstname, lastname, hashedpassword FROM users WHERE email='%s' LIMIT 1", email)

	rows, err := Conn.Query(query)
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id, &firstname, &lastname, &hashedPassword)
		if err != nil {
			log.Println(err)
		}
	}

	bPassword := []byte(password)
	bHash := []byte(hashedPassword)

	err = bcrypt.CompareHashAndPassword(bHash, bPassword)
	log.Println("PW Compare: ", err)
	if err == nil {
		ok = true
		user.Id = id
		user.FirstName = firstname
		user.LastName = lastname
		user.Email = email
	} else {
		ok = false
	}

	return user, ok
}

func AddLink(user User, link string) (result sql.Result, err error) {

	query := fmt.Sprintf("INSERT INTO links (link, userid) VALUES ('%s', '%s')", link, user.Id)

	result, err = Conn.Exec(query)
	log.Printf("Add, Result: %s, err:, %s", result, err)

	return result, err
}

func deleteLink(user User, link string) (result sql.Result, err error) {

	query := fmt.Sprintf("DELETE FROM links WHERE link='%s' AND userid='%s'", link, user.Id)

	result, err = Conn.Exec(query)
	log.Printf("Delete, Result: %s, err:, %s", result, err)

	return result, err
}

func GetUser(userEmail string) (u User, err error) {
	var (
		id        string
		email     string
		firstname string
		lastname  string
	)

	query := fmt.Sprintf("SELECT * FROM users WJHERE email='%s' LIMIT 1", userEmail)

	rows, err := Conn.Query(query)
	if err != nil {
		log.Println(err)
	}

	for rows.Next() {
		err = rows.Scan(&id, &email, &firstname, &lastname)
		if err != nil {
			log.Println(err)
		}
		u = User{id, email, firstname, lastname}
	}

	return u, err
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
		log.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&id, &email, &firstname, &lastname)
		if err != nil {
			log.Println(err)
		}
		u := User{id, email, firstname, lastname}
		allRows = append(allRows, u)
	}

	return allRows, err
}

func GetAllLinks(user User) (profile Profile, err error) {

	var (
		link string
		allLinks []string
	)

	query := fmt.Sprintf("SELECT link FROM links WHERE userId='%s'", user.Id)
	log.Println("Q: ", query)

	rows, err := Conn.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	log.Println("Checking for links")
	for rows.Next() {
		err := rows.Scan(&link)
		if err != nil {
			fmt.Print(err)
		}
		log.Println("Link: ", link)
		allLinks = append(allLinks, link)
	}

	profile.User = user
	profile.Links = allLinks

	return profile, err
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
