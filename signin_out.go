package main

import (
	"database/sql"
	"log"
	"net/http"
)

func signoutPost(w http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("mysql", dbCreds)
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

}

func signoutGet(w http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("mysql", dbCreds)
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

	_, err = stmt.Exec(boatname, hazards, damage)
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
		BoatList   []string
		IsLoggedIn bool
	}

	pageData := page{
		BoatList:   boats,
		IsLoggedIn: true,
	}

	tpl.ExecuteTemplate(w, "signin.html", pageData)

}
