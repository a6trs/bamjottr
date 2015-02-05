package grass

import (
	"../soil"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/url"
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
		if post.Author != accountInSession(w, r) {
			redirectWithError(w, r, "Don't play with others' things without permission ^^<br>Leave a comment to tell the writer to correct the mistake if there <i>is</i> one.", "/error")
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

// @url /comment/{pstid:[0-9]+} [POST]
func CommentHandler(w http.ResponseWriter, r *http.Request) {
	if accountInSession(w, r) == -1 {
		redirectWithError(w, r, "So you want to be anonymous? We'll soon support it, but I'm sorry you can't do that right now.", "/login/"+url.QueryEscape(r.URL.Path[1:]))
		return
	}
	vars := mux.Vars(r)
	pstid, _ := strconv.Atoi(vars["pstid"])
	// No need to check whether the post exists.
	// Comments with invalid post IDs won't be displayed and can be easily removed.
	text := r.FormValue("comment")
	cmt := &soil.Comment{ID: -1, PostID: pstid, Text: text, Author: accountInSession(w, r), ReplyFor: -1}
	// The key doesn't matter here. Soon it won't matter in any Save() calls.
	if err := cmt.Save(1234); err != nil {
		redirectWithError(w, r, "Cannot leave your comment there right now. Try again layer :(<br>"+err.Error(), "/error")
	} else {
		http.Redirect(w, r, "/post/"+strconv.Itoa(pstid), http.StatusFound)
	}
}
