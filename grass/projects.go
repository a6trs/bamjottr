package grass

import (
	"../soil"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/url"
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

// @url /projects
// @url /projects/{page:[0-9]+}
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

// @url /project_create
// @url /project_edit/{prjid:[0-9]+}
func ProjectEditHandler(w http.ResponseWriter, r *http.Request) {
	if accountInSession(w, r) == -1 {
		redirectWithError(w, r, "Please log in before creating a project.", "/login/"+url.QueryEscape(r.URL.Path[1:]))
		return
	}
	vars := mux.Vars(r)
	prjid_s := vars["prjid"]
	prjid := -1
	if prjid_s != "" {
		prjid, _ = strconv.Atoi(prjid_s)
	}
	if r.Method == "GET" {
		if prjid != -1 && !soil.HasMembership(prjid, accountInSession(w, r)) {
			redirectWithError(w, r, "You are not in the team of this project.", "/error")
			return
		}
		// Retrieve basic project data.
		var prj *soil.Project
		var members []*soil.ProjectMembershipData
		if prjid != -1 {
			prj = &soil.Project{ID: prjid}
			err := prj.Load(soil.KEY_Project_ID)
			if err != nil {
				prj = &soil.Project{ID: -1}
			}
			// Find all team members of this project.
			members, err = soil.AllMembers(prjid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else {
			prj = &soil.Project{ID: -1}
			members = []*soil.ProjectMembershipData{}
		}
		renderTemplate(w, r, "project_edit", map[string]interface{}{"prj": prj, "members": members})
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
		// Add the creator to the project's team if creating a new project
		if prjid == -1 {
			// The account ID 0 is the outsider colour
			soil.AddMembership(prj.ID, 0)
			soil.AddMembership(prj.ID, accountInSession(w, r))
		} else {
			// Update members' post colours.
			members, err := soil.AllMembers(prjid)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			for _, membership := range members {
				soil.UpdatePostColour(membership.ID,
					r.FormValue("postcolour-"+strconv.Itoa(membership.ID)))
			}
		}
		http.Redirect(w, r, fmt.Sprintf("/project/%d", prj.ID), http.StatusFound)
	}
}

// @url /project/{prjid:[0-9]+}
func ProjectPageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	prjid, _ := strconv.Atoi(vars["prjid"])
	prj := &soil.Project{ID: prjid}
	if err := prj.Load(soil.KEY_Project_ID); err != nil {
		// Should we use HTTP 303 (StatusSeeOther) here??
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}
	pstpage := soil.PostsForProject(prjid)
	if pstpage == nil {
		http.Error(w, "Unable to load posts TAT", http.StatusInternalServerError)
		return
	}
	// Get counts of different sight levels. The type is `map[int]int`
	allsights, cursight := soil.VisitAndCountSights("projects", prjid, accountInSession(w, r))
	renderTemplate(w, r, "project_page", map[string]interface{}{"prj": prj, "pstpage": pstpage, "allsights": allsights, "cursight": cursight})
}

// @url /invite/{prjid:[0-9]+}
// @url /invite/{prjid:[0-9]+}/{aid:[0-9]+}
func InviteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Load the project
	prjid, _ := strconv.Atoi(vars["prjid"])
	if prjid != -1 && !soil.HasMembership(prjid, accountInSession(w, r)) {
		redirectWithError(w, r, "You are not in the team of this project.", "/error")
		return
	}
	prj := &soil.Project{ID: prjid}
	if err := prj.Load(soil.KEY_Project_ID); err != nil {
		http.Redirect(w, r, "/projects", http.StatusSeeOther)
		return
	}
	// Check whether an account ID has been passed in
	aid, err := strconv.Atoi(vars["aid"])
	if err != nil {
		// Show the page
		sess, _ := sstore.Get(r, "flash")
		flash := sess.Flashes("errmsg")
		msg := ""
		iserr := true
		if flash != nil {
			msg = flash[0].(string)
			iserr = flash[1].(bool)
		}
		sess.Save(r, w)
		renderTemplate(w, r, "invite", map[string]interface{}{"prj": prj, "msg": msg, "iserr": iserr})
	} else {
		// Send an invitation to account #`aid`.
		errmsg := ""
		for {
			link := soil.InvitationLink(prjid, aid)
			if link == "" {
				errmsg = err.Error()
				break
			}
			s := fmt.Sprintf("Hey! Join us, the team of <a href='/project/%d'>%s</a>! Follow <a href='%s'>this link</a> to confirm.", prj.ID, prj.Title, link)
			if soil.SendNotification(accountInSession(w, r), aid, s) != nil {
				errmsg = err.Error()
				break
			}
			break
		}
		// Store error messages in flash data.
		iserr := errmsg != ""
		if !iserr {
			errmsg = "Invitation successfully sent."
		}
		sess, err := sstore.Get(r, "flash")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		sess.AddFlash(errmsg, "errmsg")
		sess.AddFlash(iserr, "errmsg")
		sess.Save(r, w)
		http.Redirect(w, r, "/invite/"+vars["prjid"], http.StatusFound)
	}
}

// @url /answer_invitation/{token}
func AnswerInvitationHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token, err := strconv.ParseInt(vars["token"], 36, 64)
	if err != nil {
		http.Redirect(w, r, "/notifications", http.StatusSeeOther)
		return
	}
	rsvp := soil.InvitationByToken(token)
	if rsvp != nil && rsvp.Receiver == accountInSession(w, r) {
		soil.AddMembership(rsvp.Project, rsvp.Receiver)
		http.Redirect(w, r, "/project/"+strconv.Itoa(rsvp.Project), http.StatusFound)
	} else {
		http.Redirect(w, r, "/notifications", http.StatusSeeOther)
	}
}
