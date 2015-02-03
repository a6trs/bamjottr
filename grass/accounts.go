package grass

import (
	"../soil"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
)

// @url /login          [GET, POST]
// @url /login/{return} [GET]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sess, _ := sstore.Get(r, "flash")
		flash := sess.Flashes("errmsg")
		msg := ""
		if flash != nil {
			msg = flash[0].(string)
		}
		sess.Save(r, w)
		renderTemplate(w, r, "login", map[string]interface{}{"loginmsg": msg})
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
			sess, err := sstore.Get(r, "flash")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			sess.AddFlash(errmsg, "errmsg")
			sess.Save(r, w)
			http.Redirect(w, r, r.URL.Path, http.StatusFound)
			return
		}
		sess, err := sstore.Get(r, "account-auth")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sess.Values["id"] = acc.ID
		sess.Save(r, w)
		http.Redirect(w, r, returnAddr, http.StatusFound)
	}
}

// @url /signup          [GET, POST]
// @url /signup/{return} [GET]
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		sess, _ := sstore.Get(r, "flash")
		flash := sess.Flashes("errmsg")
		msg := ""
		if flash != nil {
			msg = flash[0].(string)
		}
		sess.Save(r, w)
		renderTemplate(w, r, "signup", map[string]interface{}{"signupmsg": msg})
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
			sess, err := sstore.Get(r, "flash")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			sess.AddFlash(errmsg, "errmsg")
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
		sess.Values["id"] = acc.ID
		sess.Save(r, w)
		http.Redirect(w, r, returnAddr, http.StatusFound)
	}
}

// @url /logout
// @url /logout/{return}
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
	delete(sess.Values, "id")
	sess.Save(r, w)
	http.Redirect(w, r, returnAddr, http.StatusFound)
}

// @url /profedit
func ProfEditHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := sstore.Get(r, "account-auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	aid := sess.Values["id"]
	acc := &soil.Account{ID: aid.(int)}
	if acc.Load(soil.KEY_Account_ID) != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Method == "GET" {
		sess, _ := sstore.Get(r, "flash")
		flash := sess.Flashes("errmsg")
		msg := ""
		iserr := true
		if flash != nil {
			msg = flash[0].(string)
			iserr = flash[1].(bool)
		}
		sess.Save(r, w)
		renderTemplate(w, r, "profedit", map[string]interface{}{"aid": aid, "msg": msg, "iserr": iserr})
	} else {
		uname := r.FormValue("uname")
		email := r.FormValue("email")
		pwd := r.FormValue("pwd")
		errmsg := ""
		for {
			acc1 := &soil.Account{Name: uname}
			if acc1.Load(soil.KEY_Account_Name) == nil && acc1.ID != acc.ID {
				errmsg = "User name already exists"
				break
			}
			acc1 = &soil.Account{Email: email}
			if acc1.Load(soil.KEY_Account_Email) == nil && acc1.ID != acc.ID {
				errmsg = "E-mail already exists"
				break
			}
			break
		}
		if errmsg != "" {
			sess, err := sstore.Get(r, "flash")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			sess.AddFlash(errmsg, "errmsg")
			sess.AddFlash(true, "errmsg")
			sess.Save(r, w)
			http.Redirect(w, r, r.URL.Path, http.StatusFound)
			return
		}
		acc.Name = uname
		acc.Email = email
		if pwd != "" {
			acc.ChangePassword(pwd)
		}
		err = acc.Save(soil.KEY_Account_ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("Account updated: #", acc.ID)
		sess, err := sstore.Get(r, "flash")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sess.AddFlash("Profile successfully updated", "errmsg")
		sess.AddFlash(false, "errmsg")
		sess.Save(r, w)
		http.Redirect(w, r, r.URL.Path, http.StatusFound)
	}
}

// @url /notifications
func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	// Reset last read time of the account
	soil.UpdateLastReadTime(accountInSession(w, r))
	renderTemplate(w, r, "notifications", map[string]interface{}{})
}
