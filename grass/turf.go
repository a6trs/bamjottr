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
		Funcs(template.FuncMap{"validuser": validUser, "account": account}).
		ParseFiles("flowers/_html_head.html", "flowers/_topbar.html", "flowers/index.html", "flowers/login.html", "flowers/signup.html", "flowers/profedit.html")

func validUser(aid int) bool {
	acc := &soil.Account{ID: aid}
	err := acc.Load(soil.KEY_Account_ID)
	return (err == nil)
}

func account(aid int) *soil.Account {
	acc := &soil.Account{ID: aid}
	err := acc.Load(soil.KEY_Account_ID)
	if err == nil {
		return acc
	} else {
		return nil
	}
}

func renderTemplate(w http.ResponseWriter, title string, arg interface{}) {
	err := templates.ExecuteTemplate(w, title+".html", arg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
