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
		uname := r.FormValue("uname")
		//pwd := r.FormValue("pwd")
		sess, err := sstore.Get(r, "account-auth")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sess.Values["cookie-id"] = uname
		sess.Save(r, w)
		http.Redirect(w, r, returnAddr, http.StatusFound)
	}
}
