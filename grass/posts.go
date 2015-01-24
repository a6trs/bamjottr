package grass

import (
	"../soil"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func PostCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		renderTemplate(w, "post_create", map[string]interface{}{"aid": accountInSession(w, r)})
	} else {
		title := r.FormValue("title")
		body := r.FormValue("body")
		vars := mux.Vars(r)
		prjid, err := strconv.Atoi(vars["prjid"])
		if err != nil {
			http.Redirect(w, r, "/projects", http.StatusSeeOther)
			return
		}
		post := &soil.Post{ProjectID: prjid, Title: title, Body: body, Author: accountInSession(w, r), Priority: 2}
		if err := post.Save(soil.KEY_Post_ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/post/%d", post.ID), http.StatusFound)
	}
}
