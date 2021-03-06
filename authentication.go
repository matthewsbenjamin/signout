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

// is user will return bool when it matches a user in the database
func isUser(u string) bool {

	db, err := sql.Open("mysql", dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var r string
	err = db.QueryRow("SELECT email FROM adults WHERE email = ?", u).Scan(&r)

	if err != nil {
		return false
	} else if u == r {
		return true
	} else {
		return false
	}
}

func loginGet(w http.ResponseWriter, req *http.Request, e error) {

	tpl.ExecuteTemplate(w, "login.html", nil)

}

func loginPost(w http.ResponseWriter, req *http.Request) {

	// parse form
	req.ParseForm()

	uname := req.FormValue("uname")
	pwd := req.FormValue("pwd")
	persist := req.FormValue("persist") == "on"

	if !isUser(uname) {
		tpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	// confirmation that pwd exists
	// scan result into p
	u, err := getUserFromEmail(uname)
	if err != nil {
		// password error - redirect to /login
		tpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

	r := CheckPasswordHash(pwd, u.Pwd)
	if !r {
		// TODO - change this to loginfail - which is login with an error message
		tpl.ExecuteTemplate(w, "login.html", nil)
		return
	}

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

// HashPassword should have some documentation
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash should have some documentation
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// returns bool (true) if the BOAT ALREADY EXISTS
func validateBoatName(req *http.Request, boat string) bool {

	db, err := sql.Open("mysql", dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	u, err := getUserFromRequest(req)
	if err != nil {
		log.Fatal(err)
	}

	if u.ClubVerified && u.EmailVerified && u.Admin {
		// then it's a valid user
		// and you can test the boatname

		var b int
		r := db.QueryRow("SELECT count(1) FROM boat_locations WHERE club = ? AND boat_name = ?", u.Club)
		err = r.Scan(&b)

		// if there is no error
		if err != nil {
			return false
		}

		// if b isn't anything, then there's no boat matching that name
		if b == 0 {
			return false
		} else {
			return true
		}

	}

	return false

}
