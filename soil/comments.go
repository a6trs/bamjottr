package soil

import (
	"time"
)

type Comment struct {
	ID        int
	PostID    int
	Text      string
	Author    int
	ReplyFor  int
	CreatedAt time.Time
}

func init_Comment() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS comments (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		postid INTEGER,
		text TEXT,
		author INTEGER,
		reply_for INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(postid) REFERENCES posts(id),
		FOREIGN KEY(author) REFERENCES accounts(id),
		FOREIGN KEY(reply_for) REFERENCES comments(id)
	)`)
	return err
}

func (this *Comment) Find(key int) int {
	// Comments can only be loaded by ID.
	// XXX: Do we need to check whether the ID *does* exist?
	return this.ID
}

func (this *Comment) Load(key int) error {
	// Comments can only be loaded by ID.
	if this.ID == -1 {
		return ErrRowNotFound
	}
	row := db.QueryRow(`SELECT * FROM comments WHERE id = ?`, this.ID)
	return row.Scan(&this.ID, &this.PostID, &this.Text, &this.Author, &this.ReplyFor, &this.CreatedAt)
}

func (this *Comment) Save(key int) error {
	if this.ID == -1 {
		_, err := db.Exec(`INSERT INTO comments (postid, text, author, reply_for) VALUES (?, ?, ?, ?)`, this.PostID, this.Text, this.Author, this.ReplyFor)
		return err
	}
	_, err := db.Exec(`UPDATE comments SET postid = ?, text = ?, author = ?, reply_for = ? WHERE id = ?`, this.PostID, this.Text, this.Author, this.ReplyFor, this.ID)
	return err
}

func CommentsForPost(postid int) []*Comment {
	comments := make([]*Comment, 0)
	rows, err := db.Query(`SELECT * FROM comments WHERE postid = ?`, postid)
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		cmt := &Comment{}
		if rows.Scan(&cmt.ID, &cmt.PostID, &cmt.Text, &cmt.Author, &cmt.ReplyFor, &cmt.CreatedAt) == nil {
			comments = append(comments, cmt)
		}
	}
	return comments
}
