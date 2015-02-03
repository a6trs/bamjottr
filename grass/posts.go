package grass

import (
	"../soil"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// @url /post_create/{prjid:[0-9]+}
// @url /post_edit/{pstid:[0-9]+}
func PostEditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var prjid, pstid int
	var post *soil.Post
	var err error
	if vars["prjid"] != "" {
		prjid, _ = strconv.Atoi(vars["prjid"])
		post = &soil.Post{ProjectID: prjid, Author: accountInSession(w, r)}
		pstid = -1 // Tell the page not to display any existing data
	} else if vars["pstid"] != "" {
		pstid, _ = strconv.Atoi(vars["pstid"])
		post = &soil.Post{ID: pstid}
		if err = post.Load(soil.KEY_Post_ID); err != nil {
			http.Redirect(w, r, "/projects", http.StatusSeeOther)
			return
		}
		prjid = post.ProjectID
	} else {
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}
	if r.Method == "GET" {
		renderTemplate(w, r, "post_edit", map[string]interface{}{"prjid": prjid, "pstid": pstid})
	} else {
		title := r.FormValue("title")
		body := r.FormValue("body")
		prio, err := strconv.Atoi(r.FormValue("prio"))
		if err != nil {
			prio = soil.Post_PrioLowest
		}
		post.Title = title
		post.Body = body
		post.Priority = prio
		if err := post.Save(soil.KEY_Post_ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/post/%d", post.ID), http.StatusFound)
	}
}

// @url /post/{pstid:[0-9]+}
func PostPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pstid, _ := strconv.Atoi(vars["pstid"])
	post := &soil.Post{ID: pstid}
	if err := post.Load(soil.KEY_Post_ID); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	// Get counts of different sight levels.
	allsights, cursight := soil.VisitAndCountSights("posts", pstid, accountInSession(w, r))
	renderTemplate(w, r, "post_page", map[string]interface{}{"post": post, "allsights": allsights, "cursight": cursight})
}
