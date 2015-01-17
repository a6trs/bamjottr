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
	r.HandleFunc("/", grass.IndexHandler)
	r.HandleFunc("/login", grass.LoginHandler).Methods("GET", "POST")
	r.HandleFunc("/login/{return}", grass.LoginHandler).Methods("GET", "POST")
	r.PathPrefix("/leaves").Handler(http.FileServer(http.Dir("./flowers/")))
	http.Handle("/", r)
	soil.InitDatabase()
	fmt.Println("Now serving")
	http.ListenAndServe(":14706", nil)
}
