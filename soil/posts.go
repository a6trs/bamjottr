package soil

import (
	"fmt"
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
	Post_PrioUnsaved = 23333
)

func (this *Post) Find(key int) int {
	result := -1
	var stmt string
	switch key {
	case KEY_Post_ID:
		stmt = fmt.Sprintf(`SELECT id FROM posts WHERE id = %d`, this.ID)
	case KEY_Post_Priority:
		stmt = fmt.Sprintf(`SELECT id FROM posts WHERE priority = %d`, this.Priority)
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

func (this *Post) Load(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		return ErrRowNotFound
	}
	row := db.QueryRow(fmt.Sprintf(`SELECT * FROM posts WHERE id = %d`, this.ID))
	return row.Scan(&this.ID, &this.ProjectID, &this.Title, &this.Body, &this.Author, &this.Priority, &this.CreatedAt)
}

func (this *Post) Save(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		_, err := db.Exec(fmt.Sprintf(`INSERT INTO posts (priority) VALUES (%d)`, Post_PrioUnsaved))
		if err != nil {
			return err
		}
		prio := this.Priority
		this.Priority = Post_PrioUnsaved
		this.ID = this.Find(KEY_Post_Priority)
		this.Priority = prio
	}
	stmt := fmt.Sprintf(`UPDATE posts SET prjid = %d, title = '%s', body = '%s', author = %d, priority = %d WHERE id = %d`, this.ProjectID, this.Title, this.Body, this.Author, this.Priority, this.ID)
	_, err := db.Exec(stmt)
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
		if rows.Scan(&post.ID, &post.ProjectID, &post.Title, &post.Body, &post.Author, &post.Priority, &post.CreatedAt) == nil {
			posts = append(posts, post)
		}
	}
	return posts
}
