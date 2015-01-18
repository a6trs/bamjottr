package grass

import (
	"../soil"
	"html/template"
	"net/http"
	"github.com/gorilla/sessions"
)

var sstore = sessions.NewCookieStore([]byte("these-are-very-important-yeah"))

var templates, _ =
	template.New("IDONTKNOW").
		Funcs(template.FuncMap{"validuser": validUser, "username": userName}).
		ParseFiles("flowers/_html_head.html", "flowers/_topbar.html", "flowers/index.html", "flowers/login.html", "flowers/signup.html")

func validUser(aid int) bool {
	acc := &soil.Account{ID: aid}
	err := acc.Load(soil.KEY_Account_ID)
	return (err == nil)
}

func userName(aid int) string {
	acc := &soil.Account{ID: aid}
	err := acc.Load(soil.KEY_Account_ID)
	if err == nil {
		return acc.Name
	} else {
		return ""
	}
}

func renderTemplate(w http.ResponseWriter, title string, arg interface{}) {
	err := templates.ExecuteTemplate(w, title+".html", arg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
