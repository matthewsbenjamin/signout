package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func apiHandler(w http.ResponseWriter, req *http.Request) {

	// there should be some kind of validation here - preventing unverified users
	// from accessing api

	err := req.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	q := req.URL.Query()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// TODO case statement to run through API options
	if _, ok := q["boat"]; ok {
		apiBoatValidation(w, req)
	}

}

// api user request to
func apiUserValidation(w http.ResponseWriter, req *http.Request) {

	// err := req.ParseForm()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// // q is the parsed query from the URL as a map
	// q := req.URL.Query()

	// u, err := getUserFromEmail(q["user"][0])
	// fmt.Println(err)

	// if err != nil {

	// 	io.WriteString(w, "false\n\n")

	// } else if u.Email == q["user"][0] {

	// 	io.WriteString(w, "true\n\n")

	// }

}

// api user request to get a list of the boats in that users club
// and return that as JSON. This should be called ONCE at the loading
// of the page, and then managed by the browser
func apiBoatValidation(w http.ResponseWriter, req *http.Request) {

	// create a query for the boats for that users club
	u, err := getUserFromRequest(req)
	if err != nil {
		fmt.Fprint(w, err)
	}

	// if u.ClubVerified {
	// 	getClubBoats(u.Club)
	// }

	boats, err := getClubBoats(u.Club)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(boats)

	enc := json.NewEncoder(w)
	enc.Encode(boats)

}
