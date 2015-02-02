package grass

import (
	"../soil"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func PostEditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var prjid, pstid int
	var post *soil.Post
	var err error
	if vars["prjid"] != "" {
		prjid, err = strconv.Atoi(vars["prjid"])
		if err != nil {
			http.Redirect(w, r, "/projects", http.StatusSeeOther)
			return
		}
		post = &soil.Post{ProjectID: prjid, Author: accountInSession(w, r)}
		pstid = -1 // Tell the page not to display any existing data
	} else if vars["pstid"] != "" {
		pstid, err = strconv.Atoi(vars["pstid"])
		if err != nil {
			http.Redirect(w, r, "/projects", http.StatusSeeOther)
			return
		}
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
	// Get counts of different sight levels.
	// TODO: Merge this with code in `ProjectPageHandler` [grass/projects.go]
	allsights := soil.SightCount("sights_posts", pstid)
	sight := &soil.Sight{Account: accountInSession(w, r), Target: pstid, TableName: "sights_posts"}
	err = sight.Load(soil.KEY_Sight_AccountAndTarget)
	if err != nil {
		// No records found, create a new 'Glance' record.
		sight = &soil.Sight{Account: accountInSession(w, r), Target: pstid, TableName: "sights_posts", Level: soil.Sight_Glance}
		if err = sight.Save(soil.KEY_Sight_ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		allsights[sight.Level]++
	}
	renderTemplate(w, r, "post_page", map[string]interface{}{"post": post, "allsights": allsights, "cursight": sight.Level})
}
