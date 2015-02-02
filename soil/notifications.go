package soil

import (
	"database/sql"
	"time"
)

type Notification struct {
	ID        int
	Text      string
	Sender    int
	Receiver  int
	IsRead    bool
	CreatedAt time.Time
}

func init_Notification() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS notifications (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT,
		sender INTEGER,
		receiver INTEGER,
		isread BOOLEAN DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(sender) REFERENCES accounts(id),
		FOREIGN KEY(receiver) REFERENCES accounts(id)
	)`)
	return err
}

const (
	KEY_Notification_ID = iota
)

func (this *Notification) Find(key int) int {
	result := -1
	var row *sql.Row
	switch key {
	case KEY_Notification_ID:
		row = db.QueryRow(`SELECT id FROM notifications WHERE id = ?`, this.ID)
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

func (this *Notification) Load(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		return ErrRowNotFound
	}
	row := db.QueryRow(`SELECT * FROM notifications WHERE id = ?`, this.ID)
	return row.Scan(&this.ID, &this.Text, &this.Sender, &this.Receiver, &this.IsRead, &this.CreatedAt)
}

func (this *Notification) Save(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		// Send a new notification.
		_, err := db.Exec(`INSERT INTO notifications (text, sender, receiver) VALUES (?, ?, ?)`, this.Text, this.Sender, this.Receiver)
		return err
	} else {
		// Notifications can only be modified by the `IsRead` field.
		_, err := db.Exec(`UPDATE notifications SET isread = ? WHERE id = ?`, this.IsRead, this.ID)
		return err
	}
}

func NewNotificationsCount(account *Account) int {
	row := db.QueryRow(`SELECT COUNT(*) FROM notifications WHERE receiver = ? AND created_at > ?`, account.ID, account.LastRead)
	var n int
	if row.Scan(&n) != nil {
		return 0
	} else {
		return n
	}
}

func NotificationsFor(aid int) []*Notification {
	ret := make([]*Notification, 0)
	rows, err := db.Query(`SELECT * FROM notifications WHERE receiver = ? ORDER BY created_at DESC`, aid)
	if err != nil {
		return ret
	}
	defer rows.Close()
	for rows.Next() {
		n := &Notification{}
		if rows.Scan(&n.ID, &n.Text, &n.Sender, &n.Receiver, &n.IsRead, &n.CreatedAt) == nil {
			ret = append(ret, n)
		}
	}
	return ret
}

// XXX: Do not expose this method to anything outside the program.
func SendNotification(from, to int, text string) error {
	n := &Notification{Sender: from, Receiver: to, Text: text}
	return n.Save(KEY_Notification_ID)
}
