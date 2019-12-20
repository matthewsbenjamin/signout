package main

import (
	"fmt"
	"io"
	"net/http"
)

// coming in as /api?key=value&otherKey=otherValue
// probably as /api?user=email@example.com
func apiHandler(w http.ResponseWriter, req *http.Request) {

	// TODO add validation to this user request. Fail if validation doesn't work

	err := req.ParseForm()
	if err != nil {
		//fmt.Println(err)
	}

	q := req.URL.Query()

	// is this
	if _, ok := q["user"]; ok {
		// get
		apiUserValidation(w, req)

	}

	if _, ok := q["boat"]; ok {
		// get
		apiBoatValidation(w, req)

	}

}

// api user request to
func apiUserValidation(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	// q is the parsed query from the URL as a map
	q := req.URL.Query()

	u, err := getUserFromEmail(q["user"][0])
	fmt.Println(err)

	if err != nil {

		io.WriteString(w, "false")

	} else if u.Email == q["user"][0] {

		io.WriteString(w, "true")

	}

}

// api user request to
func apiBoatValidation(w http.ResponseWriter, req *http.Request) {

	err := req.ParseForm()
	if err != nil {
		fmt.Println(err)
	}

	u, err := getUserFromRequest(req)

	// q is the parsed query from the URL as a map
	q := req.URL.Query()

	b, err := getBoatDetails(q["boat"][0], u.Club)
	if err != nil {

		io.WriteString(w, "false")

	} else if b.Name == q["boat"][0] {

		io.WriteString(w, "true")

	}

}
