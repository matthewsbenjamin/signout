package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var tpl *template.Template
var config Configs
var cred Creds
var dbCreds string

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))

	// non-persistent store of sessions
	// var Sessions map[string]User
	// get it to read the yaml config file

	// is this the best way to do this?

	config.getConf()
	dbCreds = cred.dbCred()

}

func main() {

	// File serving
	fs := http.FileServer(http.Dir("public/"))
	http.Handle("/public/", http.StripPrefix("/public/", fs))
	img := http.FileServer(http.Dir("img/"))
	http.Handle("/img/", http.StripPrefix("/img/", img))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	// Get or config.Port request handlers
	http.HandleFunc("/", index)
	http.HandleFunc("/new-boat", newBoatHandler)
	http.HandleFunc("/new-user", newUserHandler)
	http.HandleFunc("/signout", signoutHandler)
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/hazards", hazards)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logout)

	// // API
	// http.HandleFunc("/api/", apiHandler)

	fmt.Printf("###################################\nRunning on port %s\n\n", config.Port)

	http.ListenAndServe(config.Port, nil) //
}

func ping(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "OK")

}

func index(w http.ResponseWriter, req *http.Request) {

	if !isLoggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
		return
	}

	u, err := getUserStatus(req)
	if err != nil {
		http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
		return
	}

	type Page struct {
		ErrorNotification string
		Status            string
	}

	pageData := Page{
		ErrorNotification: "",
		Status:            u,
	}

	tpl.ExecuteTemplate(w, "index.html", pageData)

}
