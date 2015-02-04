package soil

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"
)

type Project struct {
	ID          int
	Title       string
	Desc        string
	Author      int
	State       int
	TitleColour string
	BannerImg   string
	BannerType  int
	CreatedAt   time.Time
}

func init_Project() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		desc TEXT,
		author INTEGER,
		state INTEGER,
		title_clr VARCHAR(7),
		banner_img TEXT,
		banner_type INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(author) REFERENCES accounts(id)
	)`)
	return err
}

const (
	KEY_Project_ID = iota
	KEY_Project_State
)

const (
	Project_StUnsaved = iota
	Project_StPurposed
	Project_StSeeded
	Project_StRooting
	Project_StSprouts
	Project_StLawn
	Project_StWood
	Project_StJungle
	Project_StForest
)

const (
	BI_Pattern = iota
	BI_Cover
)

func StateStyles(state int) (string, string) {
	switch state {
	case Project_StUnsaved:
		return "#999999", "Unsaved"
	case Project_StPurposed:
		return "#0dcfc7", "Purposed"
	case Project_StSeeded:
		return "#9f8e0e", "Seeded"
	case Project_StRooting:
		return "#e2d904", "Rooting"
	case Project_StSprouts:
		return "#76e331", "Sprouts"
	case Project_StLawn:
		return "#01dd63", "Lawn"
	case Project_StWood:
		return "#16bc08", "Wood"
	case Project_StJungle:
		return "#049b27", "Jungle"
	case Project_StForest:
		return "#046009", "Forest"
	default:
		return "#999999", "Unknown"
	}
}

// Usage: class='banner <COBT(BT)>'
func ClassOfBannerType(bitype int) string {
	switch bitype {
	case BI_Pattern:
		return "bi-pattern"
	case BI_Cover:
		return "bi-cover"
	default:
		return ""
	}
}

func (this *Project) Find(key int) int {
	result := -1
	var row *sql.Row
	switch key {
	case KEY_Project_ID:
		row = db.QueryRow(`SELECT id FROM projects WHERE id = ?`, this.ID)
	case KEY_Project_State:
		row = db.QueryRow(`SELECT id FROM projects WHERE state = ?`, this.State)
	default:
		return -1
	}
	err := row.Scan(&result)
	if err == nil {
		return result
	} else {
		return -1
	}
}

func (this *Project) Load(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		return ErrRowNotFound
	}
	row := db.QueryRow(`SELECT * FROM projects WHERE id = ?`, this.ID)
	return row.Scan(&this.ID, &this.Title, &this.Desc, &this.Author, &this.State, &this.TitleColour, &this.BannerImg, &this.BannerType, &this.CreatedAt)
}

func (this *Project) Save(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		_, err := db.Exec(`INSERT INTO projects (state) VALUES (?)`, Project_StUnsaved)
		if err != nil {
			return err
		}
		state := this.State
		this.State = Project_StUnsaved
		this.ID = this.Find(KEY_Project_State)
		this.State = state
	}
	_, err := db.Exec(`UPDATE projects SET title = ?, desc = ?, author = ?, state = ?, title_clr = ?, banner_img = ?, banner_type = ? WHERE id = ?`, this.Title, this.Desc, this.Author, this.State, this.TitleColour, this.BannerImg, this.BannerType, this.ID)
	return err
}

func NumberOfProjects() int {
	var n int
	row := db.QueryRow(`SELECT COUNT(*) FROM projects`)
	if row.Scan(&n) == nil {
		return n
	} else {
		return -1
	}
}

func RecommendProjects(prjid int) []int {
	return Recommend(prjid, "projects", 3)
}

// ========
//  Project Team Related Section
// ========

type ProjectMembershipData struct {
	ID         int
	ProjectID  int
	AccountID  int
	PostColour string
	CreatedAt  time.Time
}

func init_ProjectMembershipData() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS projects_membership (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER,
		account_id INTEGER,
		post_colour VARCHAR(7),
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(project_id) REFERENCES projects(id),
		FOREIGN KEY(account_id) REFERENCES accounts(id)
	)`)
	return err
}

func AddMembership(prjid, aid int) error {
	_, err := db.Exec(`INSERT INTO projects_membership (project_id, account_id, post_colour) VALUES (?, ?, ?)`, prjid, aid, fmt.Sprintf("#%02x%02x%02x", rand.Int() % 128 + 128, rand.Int() % 128 + 128, rand.Int() % 128 + 128))
	return err
}

// Call with **great care**!! This operation can **not** be undone!
func RemoveMembership(prjid, aid int) error {
	_, err := db.Exec(`DELETE FROM projects_membership WHERE project_id = ? AND account_id = ?`, prjid, aid)
	return err
}

func UpdatePostColour(recid int, postColour string) error {
	_, err := db.Exec(`UPDATE projects_membership SET post_colour = ? WHERE id = ?`, postColour, recid)
	return err
}

func HasMembership(prjid, aid int) bool {
	row := db.QueryRow(`SELECT COUNT(*) FROM projects_membership WHERE project_id = ? AND account_id = ?`, prjid, aid)
	var count int
	if row.Scan(&count) != nil {
		return false
	} else {
		return (count > 0)
	}
}

func GetPostColour(prjid, aid int) string {
	row := db.QueryRow(`SELECT post_colour FROM projects_membership WHERE project_id = ? AND account_id = ?`, prjid, aid)
	var s string
	if row.Scan(&s) == nil {
		return s
	} else {
		return ""
	}
}

func AllMembers(prjid int) ([]*ProjectMembershipData, error) {
	rows, err := db.Query(`SELECT * FROM projects_membership WHERE project_id = ?`, prjid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ret := []*ProjectMembershipData{}
	for rows.Next() {
		ptm := &ProjectMembershipData{}
		if rows.Scan(&ptm.ID, &ptm.ProjectID, &ptm.AccountID, &ptm.PostColour, &ptm.CreatedAt) == nil {
			ret = append(ret, ptm)
		}
	}
	if len(ret) == 0 {
		return ret, ErrMembersNotFound
	} else {
		return ret, nil
	}
}
