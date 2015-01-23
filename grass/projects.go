package grass

import (
	"../soil"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	prjpage := make([]*soil.Project, 0)
	for i := 1; i <= 10; i++ {
		prj := &soil.Project{ID: i}
		if prj.Load(soil.KEY_Project_ID) == nil {
			prjpage = append(prjpage, prj)
		}
	}
	renderTemplate(w, "projects", map[string]interface{}{"aid": accountInSession(w, r), "prjpage": prjpage})
}

func ProjectCreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		renderTemplate(w, "project_create", map[string]interface{}{"aid": accountInSession(w, r)})
	} else {
		title := r.FormValue("title")
		desc := r.FormValue("desc")
		prj := &soil.Project{Title: title, Desc: desc, Author: accountInSession(w, r), State: soil.Project_StPurposed, BannerImg: ""}
		if err := prj.Save(soil.KEY_Project_ID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, fmt.Sprintf("/project/%d", prj.ID), http.StatusFound)
	}
}

func ProjectPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	prjid, err := strconv.Atoi(vars["prjid"])
	if err != nil {
		// Should we use HTTP 303 (StatusSeeOther) here??
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}
	prj := &soil.Project{ID: prjid}
	if err := prj.Load(soil.KEY_Project_ID); err != nil {
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}
	renderTemplate(w, "project_page", map[string]interface{}{"aid": accountInSession(w, r), "prj": prj})
}
