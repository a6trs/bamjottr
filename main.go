package main

import (
	// stackoverflow.com/q/10687627
	"./grass"
	"./soil"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
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
	r.HandleFunc("/logout", grass.LogoutHandler)
	r.HandleFunc("/logout/{return}", grass.LogoutHandler)
	r.HandleFunc("/profedit", grass.ProfEditHandler)
	// Projects-related [grass/projects.go]
	r.HandleFunc("/projects", grass.ProjectsHandler)
	r.HandleFunc("/project_create", grass.ProjectCreateHandler)
	r.HandleFunc("/project/{prjid:[0-9]+}", grass.ProjectPageHandler)
	// Static file server
	r.PathPrefix("/leaves").Handler(http.FileServer(http.Dir("./flowers/")))
	http.Handle("/", r)
	soil.InitDatabase()
	fmt.Println("Now serving")
	http.ListenAndServe(":14706", nil)
}
