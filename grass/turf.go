package grass

import (
	"../soil"
	"html/template"
	"net/http"
	"github.com/gorilla/sessions"
)

var sess = sessions.NewCookieStore([]byte("these-are-very-important-yeah"))

var templates, _ =
	template.New("IDONTKNOW").
		Funcs(template.FuncMap{"validuser": validUser, "username": userName}).
		ParseFiles("flowers/_html_head.html", "flowers/_topbar.html", "flowers/index.html", "flowers/login.html")

func validUser(cookie string) bool {
	acc := &soil.Account{Name: cookie}
	err := acc.Load(soil.KEY_Account_Name)
	return (err == nil)
}

func userName(cookie string) string {
	acc := &soil.Account{Name: cookie}
	err := acc.Load(soil.KEY_Account_Name)
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
