package main

import (
	"encoding/gob"
	"fmt"
	"github.com/gorilla/sessions"
	. "github.com/palex88/lttp/db"
	. "github.com/palex88/lttp/user"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var store *sessions.CookieStore

var templates = template.Must(template.ParseFiles(
	"edit.html",
	"view.html",
	"pages/home.html",
	"pages/account.html",
	"pages/login.html"))

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

type Page struct {
	Title string
	Body  []byte
}

func init() {

	log.Println("init")
	gob.Register(User{})

	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	store.Options = &sessions.Options{
		Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		HttpOnly: true,
	}
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	//var user User

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user := session.Values["name"]
	log.Println("Session user: ", user)

	//if len(u) > 0 {
	//	user = User{
	//		FirstName: "",
	//		Email:     u,
	//		LastName:  "",
	//	}
	//}

	t, _ := template.ParseFiles("pages/home.html")
	if (User{}) == user {
		t.Execute(w, nil)
	} else {
		t.Execute(w, user)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	name := session.Values["name"]
	if (User{}) != name {
		log.Println("User null")
		http.Redirect(w, r, "/home/", http.StatusSeeOther)
	}

	if r.Method == "GET" {
		t, _ := template.ParseFiles("pages/login.html")
		t.Execute(w, nil)
	}

	if r.Method == "POST" {
		err = r.ParseForm()
		//if err == nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//}
		email := r.Form["email"][0]
		password := r.Form["password"]
		log.Printf("E: %s, P: %s\n", email, password)

		user := User{Email: email}

		session.Values["name"] = user
		err = session.Save(r, w)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login/", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/home/", http.StatusFound)
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
		http.Redirect(w, r, "/home/", http.StatusSeeOther)
	}

	session.Values["name"] = User{}
	err = session.Save(r, w)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/home/", http.StatusFound)
}

func createAccountHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	name := session.Values["name"]
	if name != (User{}) {
		http.Redirect(w, r, "/home/", http.StatusSeeOther)
	}

	if r.Method == "GET" {
		t, _ := template.ParseFiles("pages/createuser.html")
		t.Execute(w, nil)
	} else if r.Method == "POST" {
		email := r.FormValue("email")
		firstName := r.FormValue("firstname")
		lastName := r.FormValue("lastname")
		password := r.FormValue("password")
		result, err := CreateUser(email, firstName, lastName, password)
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[0])
	}
}
