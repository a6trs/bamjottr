package grass

import (
	"../soil"
	"net/http"
	"net/url"
	"github.com/gorilla/mux"
	"fmt"
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
		sess.Values["cookie-id"] = acc.Name
		sess.Save(r, w)
		http.Redirect(w, r, returnAddr, http.StatusFound)
	}
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sess, _ := sstore.Get(r, "signup-msg")
		flash := sess.Flashes()
		msg := ""
		if flash != nil {
			msg = flash[0].(string)
		}
		sess.Save(r, w)
		renderTemplate(w, "signup", map[string]interface{}{"signupmsg": msg})
	} else {
		vars := mux.Vars(r)
		returnAddr, err := url.QueryUnescape(vars["return"])
		if err != nil {
			returnAddr = "/"
		}
		uname := r.FormValue("uname")
		email := r.FormValue("email")
		pwd := r.FormValue("pwd")
		errmsg := ""
		var acc *soil.Account
		for {
			acc = &soil.Account{Name: uname}
			if acc.Load(soil.KEY_Account_Name) == nil {
				errmsg = "User name already exists"
				break
			}
			acc = &soil.Account{Email: email}
			if acc.Load(soil.KEY_Account_Email) == nil {
				errmsg = "E-mail already exists"
				break
			}
			break
		}
		if errmsg != "" {
			sess, err := sstore.Get(r, "signup-msg")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			sess.AddFlash(errmsg)
			sess.Save(r, w)
			http.Redirect(w, r, r.URL.Path, http.StatusFound)
			return
		}
		// Create a new account
		acc = &soil.Account{ID: -1, Name: uname, Email: email, Password: []byte(pwd)}
		err = acc.Save(soil.KEY_Account_ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Account created: #", acc.ID)
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

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	returnAddr, err := url.QueryUnescape(vars["return"])
	if err != nil {
		returnAddr = "/"
	}
	sess, err := sstore.Get(r, "account-auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sess.Values["cookie-id"] = ""
	sess.Save(r, w)
	http.Redirect(w, r, returnAddr, http.StatusFound)
}
