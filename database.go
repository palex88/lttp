package main

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	uuid2 "github.com/google/uuid"
)

var Conn *sql.DB

// Initializes the database.
// When initializing it gets the configs from either system vars,
// or from a config file if vars are not set.
func init() {

	var (
		err    error
		config Config
	)

	config.Username = os.Getenv("USERNAME")
	config.Password = os.Getenv("PASSWORD")
	config.Endpoint = os.Getenv("ENDPOINT")
	config.Database = os.Getenv("DATABASE")
	config.Port = os.Getenv("PORT")

	log.Printf("U: %s, E: %s, D: %s, P: %s", config.Username, config.Endpoint, config.Database, config.Port)

	if config == (Config{}) {
		config = ParseConfigs()
	}

	conn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		config.Username,
		config.Password,
		config.Endpoint,
		config.Port,
		config.Database)

	Conn, err = sql.Open("mysql", conn)
	if err != nil {
		log.Println(err)
	}
}

// Creates user id. Uses the UUID2 package to ensure all IDs are unique.
// Returns the UUID.
func CreateUserId() (uuidStr string) {
	uuid, err := uuid2.NewUUID()
	if err != nil {
		log.Println(err)
	}

	return uuid.String()
}

// Creates a new user, given the user email is not already in use.
// Returns the results from the update and error. If error is not nil
// then a new uer has been created.
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

// Authenticates a user given an email address and password.
// Uses bcrypt package to hash the string of the password
// So that it can be compared to the stored password.
// Returns a user object and a boolean to say if the
// given credentials are valid.
func AuthUser(email string, password string) (user User, ok bool) {

	var (
		id             string
		firstname      string
		lastname       string
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

// Adds a link to a users account.
// Returns the results from the update and an error.
// If the error is nil, then the update was successful.
func AddLink(user User, link string) (result sql.Result, err error) {

	query := fmt.Sprintf("INSERT INTO links (link, userid) VALUES ('%s', '%s')", link, user.Id)

	result, err = Conn.Exec(query)
	log.Printf("Add, Result: %s, err:, %s", result, err)

	return result, err
}

// Deletes a link from the given users. Uses the user struct
// Returns results and err.
func deleteLink(user User, link string) (result sql.Result, err error) {

	query := fmt.Sprintf("DELETE FROM links WHERE link='%s' AND userid='%s'", link, user.Id)

	result, err = Conn.Exec(query)
	log.Printf("Delete, Result: %s, err:, %s", result, err)

	return result, err
}

// Gets user details from the database.
// Returns a user struct and any error
// returned form the query.
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

// Gets a list of all users from the database.
// Returns a slice of user structs for all users
// in te database.
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

// Returns all the links from a given users account.
// Takes in a user struct and returns a profile struct
// and any errors from the query.
func GetAllLinks(user User) (profile Profile, err error) {

	var (
		link     string
		id       int
		date     string
		userId   string
		allLinks []Link
	)

	query := fmt.Sprintf("SELECT id, link, date, userid FROM links WHERE userId='%s'", user.Id)
	log.Println("Q: ", query)

	rows, err := Conn.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	log.Println("Checking for links")
	for rows.Next() {
		err := rows.Scan(&id, &link, &date, &userId)
		if err != nil {
			fmt.Print(err)
		}
		l := Link{Id: id, LinkUrl: link, CreateDate: date, UserId: userId}
		log.Println("Link: ", link)
		allLinks = append(allLinks, l)
	}

	profile.User = user
	profile.Links = allLinks

	return profile, err
}

// Gets a user id given an email address.
// Returns the id as a string and any errors.
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
