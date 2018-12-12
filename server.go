package main

import (
	"encoding/gob"
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"os"
)

var store *sessions.CookieStore

var templates = template.Must(template.ParseFiles(
	"pages/home.html",
	"pages/login.html",
	"pages/account.html",
	"pages/createaccount.html",
	"pages/header.html",
	"pages/footer.html"))

type Page struct {
	Title string
	Body  []byte
}

func init() {

	gob.Register(User{})
	gob.Register(Profile{})

	log.Printf("Session key: %s", os.Getenv("SESSION_KEY"))
	//store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	store = sessions.NewCookieStore([]byte("RVPF3qQx9qK?riUgnV$r3F(a"))
	store.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	var profile Profile

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := session.Values["name"]
	original, ok := user.(User)
	if ok {
		log.Println("Session user: ", user)
	}

	if (User{}) == original {
		templates.ExecuteTemplate(w, "home", nil)
	} else {
		profile, err = GetAllLinks(original)
		log.Printf("Profile: %v", profile)
		templates.ExecuteTemplate(w, "home", profile)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//name := session.Values["name"]
	//if name == nil || name == (User{}) {
	//	log.Println("User null")
	//	http.Redirect(w, r, "/home", http.StatusSeeOther)
	//	return
	//}

	if r.Method == "GET" {
		log.Println("GET login page")
		//flash := session.Flashes()
		//t, _ := template.ParseFiles("pages/login.html")
		templates.ExecuteTemplate(w, "login", nil)
	}

	if r.Method == "POST" {
		err = r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		email := r.Form["email"][0]
		password := r.Form["password"][0]

		user, auth := AuthUser(email, password)
		if auth {
			session.Values["name"] = user
			err = session.Save(r, w)
			if err != nil {
				log.Println(err)
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/home", http.StatusFound)
			}
		} else {
			log.Printf("Login auth failed: %s\n", email)
			session.AddFlash("error", "Username or password incorrect.")
			session.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		}
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	name := fmt.Sprint(session.Values["name"])
	if len(name) == 0 {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}

	session.Values["name"] = User{}
	err = session.Save(r, w)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/home", http.StatusFound)
}

func createAccountHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	name := session.Values["name"]
	if name != (User{}) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}

	if r.Method == "GET" {
		if flash := session.Flashes(); len(flash) > 0 {
			log.Printf("Flash: %s", flash)
		}
		templates.ExecuteTemplate(w, "createaccount", nil)
	} else if r.Method == "POST" {
		email := r.FormValue("email")
		firstName := r.FormValue("firstname")
		lastName := r.FormValue("lastname")
		password := r.FormValue("password")
		confirmPassword := r.FormValue("confirmpassword")
		if password != confirmPassword {
			log.Println("Passwords don't match")
			session.AddFlash("Passwords dont match.")
			session.Save(r, w)
			http.Redirect(w, r, "/create-account", http.StatusSeeOther)
			return
		}

		result, err := CreateUser(email, firstName, lastName, password)
		if err != nil {
			log.Println(err)
			session.AddFlash("error", "Account could not be created, try a different email.")
			session.Save(r, w)
			http.Redirect(w, r, "/create-account", http.StatusSeeOther)
		} else {
			log.Printf("New account results: %s\n", result)
			session.Values["name"] = User{Email: email, FirstName: firstName, LastName: lastName}
			session.Save(r, w)
			http.Redirect(w, r, "/home", 302)
		}
	}
}

func addLinkHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	name := session.Values["name"]
	if name == (User{}) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	user := name.(User)

	link := r.FormValue("link")
	AddLink(user, link)
	http.Redirect(w, r, "/home", http.StatusFound)
}

func deleteLinkHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	name := session.Values["name"]
	if name == (User{}) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	user := name.(User)

	keys, ok := r.URL.Query()["link"]
	log.Printf("Keys: %s, OK: %t", keys, ok)

	deleteLink(user, keys[0])
	http.Redirect(w, r, "/home", http.StatusFound)
}

func accountHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	name := session.Values["name"]
	if name == (User{}) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	user := name.(User)

	templates.ExecuteTemplate(w, "account", user)
}
