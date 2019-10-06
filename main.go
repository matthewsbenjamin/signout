package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

)

var tpl *template.Template
var AUTH string = "admin:password@tcp(signout.crbqagbsfqi5.eu-west-2.rds.amazonaws.com:3306)/signout"
var port string = ":8080"


func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))

	// User just holds a users information - does not persist because it's all 
	// held in the database
	type User struct {
		Uid string
		Club string
	}

	// non-persistent store of sessions
	// var Sessions map[string]User
}


func main() {

	// File serving
	fs := http.FileServer(http.Dir("styles/"))
	http.Handle("/styles/", http.StripPrefix("/styles/", fs))
	img := http.FileServer(http.Dir("img/"))
	http.Handle("/img/", http.StripPrefix("/img/", img))
	sc := http.FileServer(http.Dir("scripts/"))
	http.Handle("/scripts/", http.StripPrefix("/scripts/", sc))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	// Get or Port request handlers
	http.HandleFunc("/", index)
	http.HandleFunc("/new-boat", newBoatHandler)
	http.HandleFunc("/new-user", newUserHandler)
	http.HandleFunc("/signout", signoutHandler)
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/hazards", hazards)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/logout", logout)

	
	fmt.Printf("###################################\nRunning on port %s\n\n", port)
	http.ListenAndServe(port, nil) // 
}

func ping(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "OK")

}

func index(w http.ResponseWriter, req *http.Request) {

	if isLoggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
	}

	tpl.ExecuteTemplate(w, "index.html", nil)

}


func logout(w http.ResponseWriter, req *http.Request) {
	
	if isLoggedIn(req) {
		// do something with the cookie - remove and redirect to index
		fmt.Println("user logged out")
	}

	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)

}

func newUserHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodGet {
		newUserGet(w, req)
	}

	if req.Method == http.MethodPost {
		newUserPost(w, req)
	}
}

func newUserGet(w http.ResponseWriter, req *http.Request) {

	tpl.ExecuteTemplate(w, "new-user.html", nil)
}

func newUserPost(w http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("mysql", AUTH)
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

	stmt, err := db.Prepare("INSERT INTO adults (email, name) VALUES (?, ?)")
	if err != nil {
		http.Error(w, "Statement preparation error", 500)
	}

	_, err = stmt.Exec(email, name)
	if err != nil {
		http.Error(w, "Statement execution error", 500)
	}

	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
}

func signoutHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		// Parse form
		// return the main page
		signoutPost(w, req)
	}

	if req.Method == http.MethodGet {
		signoutGet(w, req)
		// Return the sign in page
	}

}

func signoutPost(w http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("mysql", AUTH)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO transactions (boat_name, adult, signout) VALUES (?, ?, TRUE)")
	if err != nil {
		http.Error(w, "Statement preparation error", 500)
	}

	err = req.ParseForm()
	if err != nil {
		http.Error(w, "Login form parse error", 500)
	}

	boatname := req.FormValue("boat")
	adult := req.FormValue("adult")

	_, err = stmt.Exec(boatname, adult)
	if err != nil {
		http.Error(w, "Statement execution error", 500)
	}

	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
	// TODO - Danger page

}

func signoutGet(w http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("mysql", AUTH)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var adults []string
	var adult string
	rows, err := db.Query("SELECT name FROM adults WHERE active = 1 ORDER BY name ASC")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&adult)
		adults = append(adults, adult)
	}

	var boats []string
	var boat string
	brows, err := db.Query("SELECT boat_name FROM boat_locations WHERE on_water = 0 ORDER BY boat_name ASC")
	if err != nil {
		log.Fatal(err)
	}
	for brows.Next() {
		err = brows.Scan(&boat)
		boats = append(boats, boat)
	}

	type page struct {
		BoatList  []string
		AdultList []string
	}

	pageData := page{
		BoatList:  boats,
		AdultList: adults,
	}

	tpl.ExecuteTemplate(w, "signout.html", pageData)

}

func signinHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		signinPost(w, req)
	}

	if req.Method == http.MethodGet {
		signinGet(w, req)
	}

}

func signinPost(w http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("mysql", AUTH)
	// TODO set this up for AWS
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO transactions (boat_name, hazards, damage) VALUES (?, ?, ?)")
	if err != nil {
		http.Error(w, "Statement preparation error", 500)
	}

	err = req.ParseForm()
	if err != nil {
		http.Error(w, "Login form parse error", 500)
	}

	boatname := req.FormValue("boat")
	hazards := req.FormValue("hazards")
	damage := req.FormValue("damage")
	fmt.Println(boatname)

	_, err = stmt.Exec(boatname, hazards, damage)
	if err != nil {
		http.Error(w, "Statement execution error", 500)
	}

	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
	// TODO - Danger page

}

func signinGet(w http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("mysql", AUTH)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var boats []string
	var boat string
	brows, err := db.Query("SELECT boat_name FROM boat_locations WHERE on_water = 1 ORDER BY boat_name ASC")
	if err != nil {
		log.Fatal(err)
	}
	for brows.Next() {
		err = brows.Scan(&boat)
		boats = append(boats, boat)
	}

	type page struct {
		BoatList []string
	}

	pageData := page{
		BoatList: boats,
	}

	tpl.ExecuteTemplate(w, "signin.html", pageData)

}

func hazards(w http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("mysql", AUTH)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	type hazard struct {
		Timestamp   string
		Description string
	}

	hazards := []hazard{}
	var timestamp string
	var description string

	rows, err := db.Query("SELECT timestamp, hazards FROM transactions WHERE hazards IS NOT NULL")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&timestamp, &description)

		hazards = append(hazards, hazard{timestamp, description})
	}

	type page struct {
		HazardList []hazard
	}

	pageData := page{
		HazardList: hazards,
	}

	tpl.ExecuteTemplate(w, "hazards.html", pageData)
}

func newBoatHandler(w http.ResponseWriter, req *http.Request) {

	// Change this so that only admin people can sign in
	// Cookies etc

	if req.Method == http.MethodGet {
		newBoatGet(w, req)
	}

	if req.Method == http.MethodPost {
		newBoatPost(w, req)
	}
}

func newBoatGet(w http.ResponseWriter, req *http.Request) {

	tpl.ExecuteTemplate(w, "new-boat.html", nil)
}

func newBoatPost(w http.ResponseWriter, req *http.Request) {

	req.ParseForm()

	db, err := sql.Open("mysql", AUTH)
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


// bool whether the user is logged in or not
func isLoggedIn(req *http.Request) bool {

	// get cookie
	// if there's no sid cookie -- redirect to login page

	c, err := req.Cookie("sid")

	if err == http.ErrNoCookie {
		return false
	}

	// query the database with the UID
	// if nil - then return false
	// prevent injection somehow?

	db, err := sql.Open("mysql", AUTH)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var u string
	err = db.QueryRow("SELECT user FROM sessions WHERE sid = '?'", c).Scan(&u)

	return true
}

func loginHandler(w http.ResponseWriter, req *http.Request) {

	// if the user is logged in already - rdr index
	if isLoggedIn(req) {
		http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
	}

	if req.Method == http.MethodGet {
		loginGet(w, req, nil)
	}

	if req.Method == http.MethodPost {
		loginPost(w, req)
	}
}

func loginGet(w http.ResponseWriter, req *http.Request, e error) {

	tpl.ExecuteTemplate(w, "login.html", nil)

}

func loginPost(w http.ResponseWriter, req *http.Request) {

	// parse form
	req.ParseForm()

	uname := req.FormValue("username")
	pwd, err := hashPassword(req.FormValue("pwd"))
	if err != nil {
		http.Error(w, "Authentication error", 500)
	}
	persist := req.FormValue("persist") == "on"


	db, err := sql.Open("mysql", AUTH)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	
	// confirmation that pwd exists
	// scan result into p
	var c string
	err = db.QueryRow("SELECT pwd FROM adults WHERE email = ? AND expired", pwd).Scan(&c)
	if err == sql.ErrNoRows {
		M := "Error: Username not found"
		tpl.ExecuteTemplate(w, "login.html", M)
	}

	// get the result - r - of the password hash
	r := bcrypt.CompareHashAndPassword([]byte(pwd), []byte(c))
	
	if r == nil { // unsuccesful password

		M := errors.New("Error: Incorrect Password")
		tpl.ExecuteTemplate(w, "login.html", M)

	} else { // succesful password

		// set a cookie - sid
		id, err := uuid.NewV4()
		if err != nil {
			http.Error(w, "UUID Failed", 500)
		}
	
		cook := http.Cookie{
			Name:  "sid",
			Value: id.String(),
			Path:  "/",
		}
	
		if persist {
			cook.MaxAge = int(365 * 24 * time.Hour)
			http.SetCookie(w, &cook)
	
		} else {
			http.SetCookie(w, &cook)
	
		}

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
	// now add this to the database of sid
}


// authenticate will return true if the user is logged in
// todo - enrich the return value OR create userData func to return logged in data
func authenticate(req *http.Request) bool {

	uid, err := req.Cookie("uid")
	if err != http.ErrNoCookie {

		// they already have a cookie
		db, err := sql.Open("mysql", AUTH)
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


// compare a hash (from DB with paswrord)
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}


// hashPassword - this will hash the user's supplied password - for registration (initial storage)
// and for user logins
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}