package main

import (
	"database/sql"
	"log"
	"net/http"
)

// User contains user information. not password
type User struct {
	Email         string
	Name          string
	Pwd           string
	Club          string
	EmailVerified string
	ClubVerified  string
}

func getUserFromSID(req *http.Request) (User, error) {

	c, err := req.Cookie("sid")
	if err != nil {
		log.Println(err)
	}

	// get uid - compare w. database
	// TODO make this check a map of users? - have to purge expred sessions
	db, err := sql.Open("mysql", dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var u User
	r := db.QueryRow("SELECT email, name, club, email_verified, club_verified FROM sessions WHERE sid = '?' AND active = 1", c)
	// if err == sql.ErrNoRows {
	// 	return u, errors.New("sid does not match")
	// }

	r.Scan(&u.Email, &u.Name, &u.Club, &u.Pwd, &u.EmailVerified, &u.ClubVerified)

	return u, nil

}

func getUserFromEmail(e string) (User, error) {

	// using uid - compare w. database
	db, err := sql.Open("mysql", dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var u User
	r := db.QueryRow("SELECT email, pwd, name, club, email_verified, club_verified FROM adults WHERE email = '?' AND active = 1", e)
	// if err == sql.ErrNoRows {
	// 	return u, errors.New("sid does not match")
	// }

	r.Scan(&u.Email, &u.Pwd, &u.Name, &u.Club, &u.EmailVerified, &u.ClubVerified)

	return u, nil
}

func getClubs() []string {

	db, err := sql.Open("mysql", dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var clubs []string
	var club string
	rows, err := db.Query("SELECT DISTINCT club FROM adults ORDER BY club ASC")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err = rows.Scan(&club)
		clubs = append(clubs, club)
	}

	return clubs

}
