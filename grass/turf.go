package grass

import (
	"../soil"
	"fmt"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"regexp"
	"time"
)

var sstore = sessions.NewCookieStore([]byte("these-are-very-important-yeah"))

var templates, _ = template.New("IDONTKNOW").
	Funcs(template.FuncMap{
	"validuser": validUser, "account": account, "project": project, "post": post,
	"newNotificationsCount": soil.NewNotificationsCount, "notificationsFor": soil.NotificationsFor,
	"outsider": outsider, "outsider_colour": outsider_colour,
	"recommendedPrjs": recommendedPrjs, "recommendedPsts": recommendedPsts,
	"bannerclass": soil.ClassOfBannerType, "statebadge": stateBadge, "priobadge": priorityBadge,
	"sum": int_sum, "difference": int_difference, "product": int_product,
	"plus": int_sum, "minus": int_difference, "mul": int_product,
	"raw": rawhtml, "timestr": timestr, "sametime": sametime, "nutshell": nutshell, "autoselitem": autoSelectItem}).
	ParseFiles("stalks/_html_head.html", "stalks/_topbar.html", "stalks/_icons.svg", "stalks/_project_banner.html", "stalks/_sight_dropdown.html", "stalks/_alsolike.html", "stalks/_emojify.html", "stalks/index.html", "stalks/login.html", "stalks/signup.html", "stalks/notifications.html", "stalks/profedit.html", "stalks/projects.html", "stalks/project_edit.html", "stalks/invite.html", "stalks/project_page.html", "stalks/post_edit.html", "stalks/post_page.html")

func validUser(aid int) bool {
	acc := &soil.Account{ID: aid}
	err := acc.Load(soil.KEY_Account_ID)
	return (err == nil)
}

func account(aid int) *soil.Account {
	acc := &soil.Account{ID: aid}
	err := acc.Load(soil.KEY_Account_ID)
	if err == nil {
		return acc
	} else {
		return nil
	}
}

func project(prjid int) *soil.Project {
	prj := &soil.Project{ID: prjid}
	err := prj.Load(soil.KEY_Project_ID)
	if err == nil {
		return prj
	} else {
		return nil
	}
}

func post(pstid int) *soil.Post {
	pst := &soil.Post{ID: pstid}
	err := pst.Load(soil.KEY_Post_ID)
	if err == nil {
		return pst
	} else {
		return &soil.Post{Title: "", Body: "Write something to tell them...", Priority: soil.Post_PrioHighest}
	}
}

func outsider(project *soil.Project, aid int) bool {
	members, err := soil.GetMembers(project.ID)
	if err != nil {
		return false
	}
	for _, member := range members {
		if aid == member {
			return false
		}
	}
	return true
}

func outsider_colour(project *soil.Project) string {
	return "#ccf"
}

func recommendedPrjs(from int) []*soil.Project {
	rcmlist := soil.RecommendProjects(from)
	ret := make([]*soil.Project, len(rcmlist))
	for i, id := range rcmlist {
		ret[i] = project(id)
	}
	return ret
}

func recommendedPsts(from int) []*soil.Post {
	rcmlist := soil.RecommendPosts(from)
	ret := make([]*soil.Post, len(rcmlist))
	for i, id := range rcmlist {
		ret[i] = post(id)
	}
	return ret
}

func stateBadge(state int) template.HTML {
	bg, name := soil.StateStyles(state)
	return rawhtml(fmt.Sprintf("<span class='am-badge am-round am-text-default' style='background-color: %s'>%s</span>", bg, name))
}

func priorityBadge(prio int) template.HTML {
	var r, g, b int
	switch {
	case prio <= 0: // barrier: red
		r, g, b = 255, 0, 0
	case prio <= 1:
		r, g, b = 255, 32, 32
	case prio <= 2:
		r, g, b = 255, 64, 32
	case prio <= 3:
		r, g, b = 255, 96, 64
	case prio <= 4:
		r, g, b = 255, 108, 64
	case prio <= 5: // barrier: orange
		r, g, b = 255, 144, 96
	case prio <= 7:
		r, g, b = 255, 192, 48
	case prio <= 10: // barrier: yellow
		r, g, b = 216, 216, 0
	case prio <= 20:
		r, g, b = 108, 216, 0
	case prio <= 50: // barrier: green
		r, g, b = 0, 216, 0
	case prio <= 100:
		r, g, b = 0, 144, 255
	case prio <= 200: // barrier: blue
		r, g, b = 0, 0, 255
	case prio <= 500: // barrier: violet
		r, g, b = 255, 0, 255
	default:
		r, g, b = 168, 168, 168
	}
	return rawhtml(fmt.Sprintf("<span class='am-badge am-radius am-text-sm' style='background-color: #%02x%02x%02x'>%d</span>", r, g, b, prio))
}

func int_sum(a, b int) int {
	return a + b
}

func int_difference(a, b int) int {
	return a - b
}

func int_product(a, b int) int {
	return a * b
}

func rawhtml(s string) template.HTML {
	return template.HTML(s)
}

func timestr(t time.Time) string {
	return t.Format(time.RFC822)
}

func sametime(t1, t2 time.Time) bool {
	d := t1.Unix() - t2.Unix()
	return (d > -5) && (d < 5)
}

func nutshell(body string) string {
	// Simply remove all HTML tags.
	r, _ := regexp.Compile(`<\/?\w+(?:\s+\w+=['"].*['"])*>`)
	body = r.ReplaceAllString(body, "")
	// Help on rune arrays:
	// http://www.cnblogs.com/howDo/archive/2013/04/20/GoLang-String.html
	br := []rune(body)
	if len(body) <= 80 {
		return string(br)
	} else {
		return string(br[:80]) + "..."
	}
}

// stackoverflow.com/q/3518002
func autoSelectItem(targetval int, value int, optstr string) template.HTML {
	var selected string
	if targetval == value {
		selected = "selected='selected'"
	} else {
		selected = ""
	}
	return rawhtml(fmt.Sprintf("<option value='%d' %s>%s</option>", value, selected, optstr))
}

func renderTemplate(w http.ResponseWriter, r *http.Request, title string, arg map[string]interface{}) {
	arg["aid"] = accountInSession(w, r)
	err := templates.ExecuteTemplate(w, title+".html", arg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func accountInSession(w http.ResponseWriter, r *http.Request) int {
	sess, err := sstore.Get(r, "account-auth")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return -1
	}
	s := sess.Values["id"]
	if s == nil {
		s = -1
	}
	return s.(int)
}
