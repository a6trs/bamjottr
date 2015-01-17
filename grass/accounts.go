package grass

import (
	"../soil"
	"net/http"
	"net/url"
	"github.com/gorilla/mux"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sess, _ := sstore.Get(r, "login-msg")
		flash := sess.Flashes()
		msg := ""
		if flash != nil {
			msg = flash[0].(string)
		}
		sess.Save(r, w)
		renderTemplate(w, "login", map[string]interface{}{"loginmsg": msg})
	} else {
		vars := mux.Vars(r)
		returnAddr, err := url.QueryUnescape(vars["return"])
		if err != nil {
			returnAddr = "/"
		}
		uname := r.FormValue("uname")
		pwd := r.FormValue("pwd")
		errmsg := ""
		var acc *soil.Account
		for {
			acc = &soil.Account{Name: uname}
			if acc.Load(soil.KEY_Account_Name) == nil {
				break
			}
			acc = &soil.Account{Email: uname}
			if acc.Load(soil.KEY_Account_Email) == nil {
				break
			}
			errmsg = "The account does not exist"
			break
		}
		if errmsg == "" && !acc.MatchesPassword([]byte(pwd)) {
			errmsg = "Incorrect password"
		}
		if errmsg != "" {
			sess, err := sstore.Get(r, "login-msg")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			sess.AddFlash(errmsg)
			sess.Save(r, w)
			http.Redirect(w, r, r.URL.Path, http.StatusFound)
			return
		}
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
