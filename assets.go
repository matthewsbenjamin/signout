package main

import (
	"database/sql"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func newBoatGet(w http.ResponseWriter, req *http.Request) {

	tpl.ExecuteTemplate(w, "new-boat.html", nil)

}

func newBoatPost(w http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	db, err := sql.Open("mysql", dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO boat_locations (boat_name) VALUES (?);")
	if err != nil {
		log.Fatal(err)
	}

	b := req.FormValue("boatName")

	_, err = stmt.Exec(b)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)

}

func newUserGet(w http.ResponseWriter, req *http.Request) {

	tpl.ExecuteTemplate(w, "new-user.html", nil)

}

func newUserPost(w http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("mysql", dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = req.ParseForm()
	if err != nil {
		http.Error(w, "Login form parse error", 500)
	}

	email := req.FormValue("email")
	name := req.FormValue("name")
	pwd := req.FormValue("pwd")
	pwdConf := req.FormValue("pwd2")
	club := req.FormValue("club")

	if pwd != pwdConf {
		// do something about it
	}

	pwdH, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "password hashing error", 500)
	}

	stmt, err := db.Prepare("INSERT INTO adults (email, name, pwd, club) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Statement preparation error", 500)
	}

	_, err = stmt.Exec(email, name, pwdH, club)
	if err != nil {
		http.Error(w, "Statement execution error", 500)
	}

	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}
