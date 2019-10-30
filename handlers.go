package main

import "net/http"

func newUserHandler(w http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodGet {
		newUserGet(w, req)
	}

	if req.Method == http.MethodPost && isLoggedIn(req) {
		newUserPost(w, req)
	} else {
		http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
	}
}

func signoutHandler(w http.ResponseWriter, req *http.Request) {

	if !isLoggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
	}

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

func signinHandler(w http.ResponseWriter, req *http.Request) {

	if !isLoggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
	}

	if req.Method == http.MethodPost {
		signinPost(w, req)
	}

	if req.Method == http.MethodGet {
		signinGet(w, req)
	}

}

func newBoatHandler(w http.ResponseWriter, req *http.Request) {

	if !isLoggedIn(req) {
		http.Redirect(w, req, "/login", http.StatusTemporaryRedirect)
	}

	// Change this so that only admin people can sign in
	// Cookies etc

	if req.Method == http.MethodGet {
		newBoatGet(w, req)
	}

	if req.Method == http.MethodPost {
		newBoatPost(w, req)
	}
}

func loginHandler(w http.ResponseWriter, req *http.Request) {

	// // if the user is logged in already - rdr index
	// if isLoggedIn(req) {
	// 	http.Redirect(w, req, "/", http.StatusTemporaryRedirect)
	// }

	if req.Method == http.MethodGet {
		loginGet(w, req, nil)
	}

	if req.Method == http.MethodPost {
		loginPost(w, req)
	}
}
