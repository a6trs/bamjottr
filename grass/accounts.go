package grass

import (
	_ "../soil"
	"net/http"
	"net/url"
	"github.com/gorilla/mux"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		renderTemplate(w, "login", nil)
	} else {
		vars := mux.Vars(r)
		returnAddr, err := url.QueryUnescape(vars["return"])
		if err != nil {
			returnAddr = "/"
		}
		//uname := r.FormValue("uname")
		//pwd := r.FormValue("pwd")
		http.Redirect(w, r, returnAddr, http.StatusFound)
	}
}
