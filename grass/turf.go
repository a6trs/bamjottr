package grass

import (
	"../soil"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"time"
)

var sstore = sessions.NewCookieStore([]byte("these-are-very-important-yeah"))

var templates, _ = template.New("IDONTKNOW").
	Funcs(template.FuncMap{"validuser": validUser, "account": account, "timestr": timestr}).
	ParseFiles("flowers/_html_head.html", "flowers/_topbar.html", "flowers/index.html", "flowers/login.html", "flowers/signup.html", "flowers/profedit.html", "flowers/projects.html", "flowers/project_create.html", "flowers/project_page.html")

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

func timestr(t time.Time) string {
	return t.Format(time.RFC822)
}

func renderTemplate(w http.ResponseWriter, title string, arg interface{}) {
	err := templates.ExecuteTemplate(w, title+".html", arg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func accountInSession(w http.ResponseWriter, r *http.Request) int {
	sess, err := sstore.Get(r, "account-auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return -1
	}
	s := sess.Values["id"]
	if s == nil {
		s = -1
	}
	return s.(int)
}
