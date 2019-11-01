package main

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// bool whether the user is logged in or not
func isLoggedIn(req *http.Request) bool {

	// get cookie
	// if there's no sid cookie -- redirect to login page

	c, err := req.Cookie("sid")
	if err == http.ErrNoCookie {
		return false
	}
	sid := strings.Split(c.String(), "=")[1]

	// query the database with the UID
	// if nil - then return false
	// prevent injection somehow?

	db, err := sql.Open("mysql", dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var u string
	err = db.QueryRow("SELECT user FROM sessions WHERE sid = ?", sid).Scan(&u)
	if err == sql.ErrNoRows {
		return false
	} else if err != nil {
		return false
	}

	return true
}

// authenticate will return true if the user is logged in
// todo - enrich the return value OR create userData func to return logged in data
// why is this needed - different from isLoggedIn?
func authenticate(req *http.Request) bool {

	uid, err := req.Cookie("uid")
	if err != http.ErrNoCookie {

		// they already have a cookie
		db, err := sql.Open("mysql", dbCreds)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// get the corresponding uid from the database
		// then scan the count of uid into ct
		var ct int
		err = db.QueryRow("SELECT count(uid) FROM logins WHERE uid =  ?", uid.Value).Scan(&ct)
		if err == sql.ErrNoRows {
			return false
		}
		if err != nil {
			log.Fatal(err)
		}

		return ct > 0
	}
	return false
}

func loginGet(w http.ResponseWriter, req *http.Request, e error) {

	type Page struct {
		LoggedIn bool
	}

	pageData := Page{
		LoggedIn: false,
	}

	tpl.ExecuteTemplate(w, "login.html", pageData)

}

func loginPost(w http.ResponseWriter, req *http.Request) {

	// parse form
	req.ParseForm()

	uname := req.FormValue("uname")
	pwd := req.FormValue("pwd")
	persist := req.FormValue("persist") == "on"

	// confirmation that pwd exists
	// scan result into p
	// TODO change this for getUserFromEmail()
	c, err := getUserFromEmail(uname)
	if err != nil {
		// password error - redirect to /login
		http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
	}

	err = bcrypt.CompareHashAndPassword([]byte(pwd), []byte(c.Pwd))
	if err == bcrypt.ErrMismatchedHashAndPassword {

		// what do I do with loginfail?
		// redirect to loginfail, collect some hidden form data and then
		// redirect to /login
		http.Redirect(w, req, "/loginfail", http.StatusTemporaryRedirect)
	} else { // succesful password

		// set a cookie - sid
		i, err := uuid.NewV4()
		if err != nil {
			http.Error(w, "UUID Failed", 500)
		}
		id := i.String()

		cook := http.Cookie{
			Name:  "sid",
			Value: id,
			Path:  "/",
		}

		if persist { // if the user has selected remember me - persistent cookie
			cook.MaxAge = int(365 * 24 * time.Hour)
			http.SetCookie(w, &cook)

		} else { // forget me - session cookie
			http.SetCookie(w, &cook)

		}

		db, err := sql.Open("mysql", dbCreds)
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// add this succesful login to the session database
		stmt, err := db.Prepare("INSERT INTO sessions (sid, user) VALUES (?, ?)")
		if err != nil {
			http.Error(w, "Statement preparation failed", 500)
		}
		// Write to database
		_, err = stmt.Exec(id, uname)
		if err != nil {
			http.Error(w, "Statement execution failed", 500)
		}
	}

	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}

func logout(w http.ResponseWriter, req *http.Request) {

	// do something with the cookie - remove and redirect to index
	c, err := req.Cookie("sid")

	if err != nil {
		http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
		return
	}

	c.MaxAge = -1
	http.SetCookie(w, c)

	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)

}
