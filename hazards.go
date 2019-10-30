package main

import (
	"database/sql"
	"log"
	"net/http"
)

func hazards(w http.ResponseWriter, req *http.Request) {

	if !isLoggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
	}

	db, err := sql.Open("mysql", dbCreds)
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
