package grass

import (
	"html/template"
	"net/http"
)

var templates = template.Must(template.ParseFiles("flowers/_html_head.html", "flowers/_topbar.html", "flowers/index.html", "flowers/login.html"))

func renderTemplate(w http.ResponseWriter, title string, arg interface{}) {
	err := templates.ExecuteTemplate(w, title+".html", arg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
