package main

import (
	// stackoverflow.com/q/10687627
	"./grass"
	"./soil"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	// Index page [grass/index.go]
	r.HandleFunc("/", grass.IndexHandler)
	// Accounts-related [grass/accounts.go]
	r.HandleFunc("/login", grass.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/login/{return}", grass.LoginHandler).Methods("GET")
	r.HandleFunc("/signup", grass.SignupHandler).Methods("GET", "POST")
	r.HandleFunc("/signup/{return}", grass.SignupHandler).Methods("GET")
	// Static file server
	r.PathPrefix("/leaves").Handler(http.FileServer(http.Dir("./flowers/")))
	http.Handle("/", r)
	soil.InitDatabase()
	fmt.Println("Now serving")
	http.ListenAndServe(":14706", nil)
}
