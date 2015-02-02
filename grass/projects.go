package grass

import (
	"../soil"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const PROJECTS_PER_PAGE = 10

func navigationDisplay(curpage, pagecnt int) []int {
	p := []int{1}    // 1 always stays here
	if 6 > curpage { // 6 = 1 + 2(omitted) + 2(near)
		for i := 2; i <= curpage; i++ {
			p = append(p, i)
		}
	} else {
		p = append(p, -1)
		for i := curpage - 2; i <= curpage; i++ {
			p = append(p, i)
		}
	}
	if pagecnt-5 < curpage {
		for i := curpage + 1; i <= pagecnt; i++ {
			p = append(p, i)
		}
	} else {
		for i := curpage + 1; i <= curpage+2; i++ {
			p = append(p, i)
		}
		p = append(p, -1)
		p = append(p, pagecnt)
	}
	return p
}

func ProjectsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	curpage, err := strconv.Atoi(vars["page"])
	if err != nil {
		curpage = 1
	}
	prjpage := make([]*soil.Project, 0)
	for i := curpage*PROJECTS_PER_PAGE - 9; i <= curpage*PROJECTS_PER_PAGE; i++ {
		prj := &soil.Project{ID: i}
		if prj.Load(soil.KEY_Project_ID) == nil {
			prjpage = append(prjpage, prj)
		}
	}
	prjcnt := soil.NumberOfProjects()
	pagecnt := 0
	// XXX: Is there a better way to do this?
	if prjcnt != -1 {
		if prjcnt%PROJECTS_PER_PAGE == 0 {
			pagecnt = prjcnt / PROJECTS_PER_PAGE
		} else {
			pagecnt = prjcnt/PROJECTS_PER_PAGE + 1
		}
	}
	navpages := navigationDisplay(curpage, pagecnt)
	renderTemplate(w, r, "projects", map[string]interface{}{"prjpage": prjpage, "curpage": curpage, "pagecnt": pagecnt, "navpages": navpages})
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
			// Build a file name with timestamp in order to prevent conflicts.
			bannerimg_file = strconv.FormatInt(time.Now().Unix(), 36) + "-" + handler.Filename
			defer file.Close()
			newfile, err := os.OpenFile("./uploads/banner_img/"+bannerimg_file,
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
		state, err := strconv.Atoi(r.FormValue("state"))
		if err != nil {
			state = 1
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
		prj = &soil.Project{ID: prjid, Title: title, Desc: desc, Author: accountInSession(w, r), State: state, TitleColour: titlecolour, BannerImg: bannerimg_file, BannerType: bannertype}
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
	// Get counts of different sight levels. The type is `map[int]int`
	allsights, cursight := soil.VisitAndCountSights("projects", prjid, accountInSession(w, r))
	renderTemplate(w, r, "project_page", map[string]interface{}{"prj": prj, "pstpage": pstpage, "allsights": allsights, "cursight": cursight})
}
