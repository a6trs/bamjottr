package soil

import (
	"database/sql"
	"time"
)

type Post struct {
	ID        int
	ProjectID int
	Title     string
	Body      string
	Author    int
	Priority  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func init_Post() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		prjid INTEGER,
		title TEXT,
		body TEXT,
		author INTEGER,
		priority INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(prjid) REFERENCES projects(id),
		FOREIGN KEY(author) REFERENCES accounts(id)
	)`)
	return err
}

const (
	KEY_Post_ID = iota
	KEY_Post_Priority
)

const (
	Post_PrioHighest = 0
	Post_PrioLowest  = 23332
	Post_PrioUnsaved = 23333
)

func (this *Post) Find(key int) int {
	result := -1
	var row *sql.Row
	switch key {
	case KEY_Post_ID:
		row = db.QueryRow(`SELECT id FROM posts WHERE id = ?`, this.ID)
	case KEY_Post_Priority:
		row = db.QueryRow(`SELECT id FROM posts WHERE priority = ?`, this.Priority)
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

func (this *Post) Load(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		return ErrRowNotFound
	}
	row := db.QueryRow(`SELECT * FROM posts WHERE id = ?`, this.ID)
	return row.Scan(&this.ID, &this.ProjectID, &this.Title, &this.Body, &this.Author, &this.Priority, &this.CreatedAt, &this.UpdatedAt)
}

func (this *Post) Save(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		_, err := db.Exec(`INSERT INTO posts (priority) VALUES (?)`, Post_PrioUnsaved)
		if err != nil {
			return err
		}
		prio := this.Priority
		this.Priority = Post_PrioUnsaved
		this.ID = this.Find(KEY_Post_Priority)
		this.Priority = prio
	}
	_, err := db.Exec(`UPDATE posts SET prjid = ?, title = ?, body = ?, author = ?, priority = ?, updated_at = datetime('now') WHERE id = ?`, this.ProjectID, this.Title, this.Body, this.Author, this.Priority, this.ID)
	return err
}

func PostsForProject(prjid int) []*Post {
	posts := make([]*Post, 0)
	rows, err := db.Query(`SELECT * FROM posts WHERE prjid = ?`, prjid)
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		post := &Post{}
		if rows.Scan(&post.ID, &post.ProjectID, &post.Title, &post.Body, &post.Author, &post.Priority, &post.CreatedAt, &post.UpdatedAt) == nil {
			posts = append(posts, post)
		}
	}
	return posts
}

func RecommendPosts(pstid int) []int {
	return Recommend(pstid, "posts", 5)
}
