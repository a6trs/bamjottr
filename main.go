package main

import (
	"./grass"
	"./soil"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", grass.IndexHandler)
	http.Handle("/", r)
	soil.InitDatabase()
	fmt.Println("Now serving")
	http.ListenAndServe(":14706", nil)
}
