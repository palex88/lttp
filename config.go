package main

import (
	"database/sql"
	"fmt"
)

var (
	Database     *sql.DB
	//OAuthConfig  *oauth2.Config
	//SessionStore sessions.Store
)

func init() {
	//OAuthConfig = configureOAuthClient("clientid", "clientsecret")

	//cookieStore := sessions.NewCookieStore([]byte("something-very-secret"))
	//cookieStore.Options = &sessions.Options{
	//	HttpOnly: true,
	//}
	//SessionStore = cookieStore
}

func configureDatabase(config Config) (db *sql.DB) {

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

//func configureOAuthClient(clientID, clientSecret string) *oauth2.Config {
//	redirectURL := os.Getenv("OAUTH2_CALLBACK")
//	if redirectURL == "" {
//		redirectURL = "http://localhost:8080/oauth2callback"
//	}
//	return &oauth2.Config{
//		ClientID:     clientID,
//		ClientSecret: clientSecret,
//		RedirectURL:  redirectURL,
//		Scopes:       []string{"email", "profile"},
//		Endpoint:     google.Endpoint,
//	}
//}
