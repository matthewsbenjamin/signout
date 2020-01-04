package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func newBoatGet(w http.ResponseWriter, req *http.Request) {

	tpl.ExecuteTemplate(w, "new-boat.html", nil)

}

func newBoatPost(w http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	b := req.FormValue("boatName")

	if validateBoatName(req, b) {

		// then the boat already exists
		// serve the page with the user status
		status, err := getUserStatus(req)
		if err != nil {
			http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
		}

		type Page struct {
			ErrorNotification string
			Status            string
		}

		pageData := Page{
			ErrorNotification: fmt.Sprintf("A boat with the name %s already exists", b),
			Status:            status,
		}

		tpl.ExecuteTemplate(w, "new-boat.html", pageData)
		return
	}

	db, err := sql.Open("mysql", dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO boat_locations (boat_name, club) VALUES (?, ?)")
	if err != nil {
		log.Fatal(err)
	}

	u, err := getUserFromRequest(req)
	if err != nil {
		log.Fatal(err)
	}

	_, err = stmt.Exec(b, u.Club)
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
