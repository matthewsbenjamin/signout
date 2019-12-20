package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// User contains user information. not password
type User struct {
	Email         string
	Name          string
	Pwd           string
	Club          string
	EmailVerified bool
	ClubVerified  bool
}

func getUserFromRequest(req *http.Request) (User, error) {

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
	r := db.QueryRow("SELECT `email`, `name`, `club`, `pwd`, `email_verified`, `club_verified` FROM sessions INNER JOIN adults ON sessions.user = adults.email WHERE sessions.sid = ?", c.Value)
	err = r.Scan(&u.Email, &u.Name, &u.Club, &u.Pwd, &u.EmailVerified, &u.ClubVerified)

	if err != nil {
		return u, err
	}

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
	r := db.QueryRow("SELECT email, pwd, name, club, email_verified, club_verified FROM adults WHERE email = ?", e)

	err = r.Scan(&u.Email, &u.Pwd, &u.Name, &u.Club, &u.EmailVerified, &u.ClubVerified)

	if err == sql.ErrNoRows {
		return u, err
	}

	if err != nil {
		return u, err
	}

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

// Boat is a boat
// note the issue with pulling the 'current' status from the db - easily overridden
type Boat struct {
	Name string
	Club string
}

func getBoatDetails(b, c string) (Boat, error) {

	db, err := sql.Open("mysql", dbCreds)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := db.QueryRow("SELECT boat_name, club FROM boat_locations WHERE boat_name = ? AND club = ?", b, c)

	var B Boat
	err = r.Scan(&B.Name, &B.Club)
	if err != nil {
		return B, err
	}

	return B, nil
}
