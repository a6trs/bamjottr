package grass

import (
	"../soil"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
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
	renderTemplate(w, r, "projects", map[string]interface{}{"prjpage": prjpage})
}

func ProjectEditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	prjid_s := vars["prjid"]
	prjid := -1
	if prjid_s != "" {
		var err error
		prjid, err = strconv.Atoi(prjid_s)
		if err != nil {
			prjid = -1
		}
	}
	if r.Method == "GET" {
		var prj *soil.Project
		if prjid != -1 {
			prj = &soil.Project{ID: prjid}
			err := prj.Load(soil.KEY_Project_ID)
			if err != nil {
				prj = &soil.Project{ID: -1}
			}
		} else {
			prj = &soil.Project{ID: -1}
		}
		renderTemplate(w, r, "project_edit", map[string]interface{}{"prj": prj})
	} else {
		r.ParseMultipartForm(16 << 20) // 16 MB of memory
		file, handler, err := r.FormFile("bannerimg")
		var bannerimg_file string
		if err == http.ErrMissingFile {
			bannerimg_file = ""
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			bannerimg_file = handler.Filename
			defer file.Close()
			newfile, err := os.OpenFile("./uploads/banner_img/"+handler.Filename,
				os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer newfile.Close()
			io.Copy(newfile, file)
		}
		title := r.FormValue("title")
		desc := r.FormValue("desc")
		titlecolour := r.FormValue("titlecolour")
		bannertype, err := strconv.Atoi(r.FormValue("bannertype"))
		if err != nil {
			bannertype = 0
		}
		// If banner image is not changed, we read it and let it remain the same.
		prj := &soil.Project{ID: prjid}
		if bannerimg_file == "" && prjid != -1 {
			err = prj.Load(soil.KEY_Project_ID)
			if err == nil {
				bannerimg_file = prj.BannerImg
			}
		}
		// If creating project, prjid would be -1 and a new row would be created.
		prj = &soil.Project{ID: prjid, Title: title, Desc: desc, Author: accountInSession(w, r), State: soil.Project_StPurposed, TitleColour: titlecolour, BannerImg: bannerimg_file, BannerType: bannertype}
		err = prj.Save(soil.KEY_Project_ID)
		if err != nil {
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
	pstpage := soil.PostsForProject(prjid)
	if pstpage == nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	renderTemplate(w, r, "project_page", map[string]interface{}{"prj": prj, "pstpage": pstpage})
}
