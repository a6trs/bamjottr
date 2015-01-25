package soil

import (
	"database/sql"
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

func StateName(state int) string {
	switch state {
	case Project_StUnsaved:
		return "Unsaved"
	case Project_StPurposed:
		return "Purposed"
	case Project_StSeeded:
		return "Seeded"
	case Project_StRooting:
		return "Rooting"
	case Project_StSprouts:
		return "Sprouts"
	case Project_StLawn:
		return "Lawn"
	case Project_StWood:
		return "Wood"
	case Project_StJungle:
		return "Jungle"
	case Project_StForest:
		return "Forest"
	default:
		return "Unknown"
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
