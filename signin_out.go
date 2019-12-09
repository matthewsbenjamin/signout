package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

func signoutPost(w http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("mysql", dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO transactions (boat_name, adult, signout, club) VALUES (?, ?, TRUE, ?)")
	if err != nil {
		http.Error(w, "Statement preparation error", 500)
	}

	err = req.ParseForm()
	if err != nil {
		http.Error(w, "Login form parse error", 500)
	}

	u, err := getUserFromSID(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("User is: %v\n", u)

	boatname := req.FormValue("boat")
	adult := req.FormValue("adult")

	_, err = stmt.Exec(boatname, adult, u.Club)
	if err != nil {
		http.Error(w, "Statement execution error", 500)
	}

	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)

}

func signoutGet(w http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("mysql", dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	u, err := getUserFromSID(req)
	if err != nil {
		http.Error(w, "Signout get user from si failed", 401)
	}

	var adults []string
	var adult string
	rows, err := db.Query("SELECT name FROM adults WHERE active = 1 AND club = ? ORDER BY name ASC", u.Club)
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&adult)
		adults = append(adults, adult)
	}

	var boats []string
	var boat string
	brows, err := db.Query("SELECT boat_name FROM boat_locations WHERE on_water = 0 AND club = ? ORDER BY boat_name ASC", u.Club)
	if err != nil {
		log.Fatal(err)
	}
	for brows.Next() {
		err = brows.Scan(&boat)
		boats = append(boats, boat)
	}

	type page struct {
		BoatList   []string
		AdultList  []string
		IsLoggedIn bool
	}

	pageData := page{
		BoatList:   boats,
		AdultList:  adults,
		IsLoggedIn: true,
	}

	tpl.ExecuteTemplate(w, "signout.html", pageData)

}

func signinPost(w http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("mysql", dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// TODO add the users club - use getUserFromSID()
	stmt, err := db.Prepare("INSERT INTO transactions (boat_name, hazards, damage, club) VALUES (?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Statement preparation error", 500)
	}

	err = req.ParseForm()
	if err != nil {
		http.Error(w, "Login form parse error", 500)
	}

	u, err := getUserFromSID(req)
	if err != nil {
		log.Fatal(err)
	}

	boatname := req.FormValue("boat")
	hazards := req.FormValue("hazards")
	damage := req.FormValue("damage")

	_, err = stmt.Exec(boatname, hazards, damage, u.Club)
	if err != nil {
		http.Error(w, "Statement execution error", 500)
	}

	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
	// TODO - Danger page

}

func signinGet(w http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("mysql", dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	u, err := getUserFromSID(req)
	if err != nil {
		http.Error(w, "Signin get user from sid failed", 401)
	}

	var boats []string
	var boat string
	brows, err := db.Query("SELECT boat_name FROM boat_locations WHERE on_water = 1 AND club = ? ORDER BY boat_name ASC", u.Club)
	if err != nil {
		log.Fatal(err)
	}
	for brows.Next() {
		err = brows.Scan(&boat)
		boats = append(boats, boat)
	}

	type page struct {
		BoatList   []string
		IsLoggedIn bool
	}

	pageData := page{
		BoatList:   boats,
		IsLoggedIn: true,
	}

	tpl.ExecuteTemplate(w, "signin.html", pageData)

}
