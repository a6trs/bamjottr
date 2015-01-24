package soil

import (
	"fmt"
	"time"
)

type Project struct {
	ID        int
	Title     string
	Desc      string
	Author    int
	State     int
	BannerImg string
	CreatedAt time.Time
}

func init_Project() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		desc TEXT,
		author INTEGER,
		state INTEGER,
		banner_img TEXT,
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
)

func (this *Project) Find(key int) int {
	result := -1
	var stmt string
	switch key {
	case KEY_Project_ID:
		stmt = fmt.Sprintf(`SELECT id FROM projects WHERE id = %d`, this.ID)
	case KEY_Project_State:
		stmt = fmt.Sprintf(`SELECT id FROM projects WHERE state = %d`, this.State)
	default:
		return -1
	}
	row := db.QueryRow(stmt)
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
	row := db.QueryRow(fmt.Sprintf(`SELECT * FROM projects WHERE id = %d`, this.ID))
	return row.Scan(&this.ID, &this.Title, &this.Desc, &this.Author, &this.State, &this.BannerImg, &this.CreatedAt)
}

func (this *Project) Save(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		_, err := db.Exec(fmt.Sprintf(`INSERT INTO projects (state) VALUES (%d)`, Project_StUnsaved))
		if err != nil {
			return err
		}
		state := this.State
		this.State = Project_StUnsaved
		this.ID = this.Find(KEY_Project_State)
		this.State = state
	}
	stmt := fmt.Sprintf(`UPDATE projects SET title = '%s', desc = '%s', author = %d, state = %d, banner_img = '%s' WHERE id = %d`, this.Title, this.Desc, this.Author, this.State, this.BannerImg, this.ID)
	_, err := db.Exec(stmt)
	return err
}
