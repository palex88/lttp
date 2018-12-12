package main

import (
	"log"
	"net/http"
	"os"
)

func main() {

	var port string

	port = os.Getenv("PORT")
	log.Println("PORT: ", port)
	if len(port) == 0 {
		port = ":8080"
	}

	http.Handle("/css/", http.FileServer(http.Dir("")))
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/home/", homeHandler)
	http.HandleFunc("/login/", loginHandler)
	http.HandleFunc("/logout/", logoutHandler)
	http.HandleFunc("/create-account/", createAccountHandler)
	http.HandleFunc("/addlink/", addLinkHandler)
	http.HandleFunc("/deletelink/", deleteLinkHandler)
	http.HandleFunc("/account/", accountHandler)

	log.Printf("Opening server on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
