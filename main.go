package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	//"time"

	_ "github.com/go-sql-driver/mysql"
)

var tpl *template.Template
var AUTH string = "root:password@tcp(localhost:3306)/signout"

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func main() {

	// File serving
	fs := http.FileServer(http.Dir("styles/"))
	http.Handle("/styles/", http.StripPrefix("/styles/", fs))
	img := http.FileServer(http.Dir("img/"))
	http.Handle("/img/", http.StripPrefix("/img/", img))
	http.Handle("/favicon.ico", http.NotFoundHandler())

	// Get or Port request handlers
	http.HandleFunc("/", index)
	// http.HandleFunc("/new-boat", newBoatHandler)
	http.HandleFunc("/new-user", newUserHandler)
	http.HandleFunc("/signout", signoutHandler)
	http.HandleFunc("/signin", signinHandler)
	http.HandleFunc("/hazards", hazards)

	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, req *http.Request) {

	tpl.ExecuteTemplate(w, "index.html", nil)

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

	// TODO handle errors if the user has already registered

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

// func newBoatHandler(w http.ResponseWriter, req *http.Request) {

// 	if req.Method == http.MethodPost {
// 		// Parse form
// 		// return the main page
// 		boatPost(w, req)
// 	}

// 	if req.Method == http.MethodGet {
// 		boatGet(w, req)
// 		// Return the sign in page
// 	}
// }

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
	// TODO set this up for AWS
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
		Timestamp string
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
