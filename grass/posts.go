package grass

import (
	"../soil"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func PostCreateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	prjid, err := strconv.Atoi(vars["prjid"])
	if err != nil {
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		renderTemplate(w, r, "post_create", map[string]interface{}{"prjid": prjid})
	} else {
		title := r.FormValue("title")
		body := r.FormValue("body")
		post := &soil.Post{ProjectID: prjid, Title: title, Body: body, Author: accountInSession(w, r), Priority: 2}
		if err := post.Save(soil.KEY_Post_ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/post/%d", post.ID), http.StatusFound)
	}
}

func PostPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pstid, err := strconv.Atoi(vars["pstid"])
	if err != nil {
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}
	post := &soil.Post{ID: pstid}
	if err := post.Load(soil.KEY_Post_ID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	renderTemplate(w, r, "post_page", map[string]interface{}{"post": post})
}
