package soil

import (
	"database/sql"
	"time"
)

type Sight struct {
	ID        int
	Account   int
	Target    int
	Level     int
	UpdatedAt time.Time
	// Stores which table it belongs to.
	TableName string
}

func init_Sight(targetTable string) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS sights_` + targetTable + ` (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		account INTEGER,
		target INTEGER,
		level INTEGER,
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(account) REFERENCES accounts(id)
		FOREIGN KEY(target) REFERENCES ` + targetTable + `(id)
	)`)
	return err
}

const (
	KEY_Sight_ID = iota
	KEY_Sight_Account
	KEY_Sight_Target
	KEY_Sight_AccountAndTarget
	KEY_Sight_Level // Used when creating new records
)

const (
	Sight_Glance = iota
	Sight_Watch
	Sight_Stare
	Sight_Unsaved = 12138
)

func (this *Sight) Find(key int) int {
	result := -1
	var row *sql.Row
	switch key {
	case KEY_Sight_ID:
		row = db.QueryRow(`SELECT id FROM `+this.TableName+` WHERE id = ?`, this.ID)
	case KEY_Sight_Account:
		row = db.QueryRow(`SELECT id FROM `+this.TableName+` WHERE account = ?`, this.Account)
	case KEY_Sight_Target:
		row = db.QueryRow(`SELECT id FROM `+this.TableName+` WHERE target = ?`, this.Target)
	case KEY_Sight_AccountAndTarget:
		row = db.QueryRow(`SELECT id FROM `+this.TableName+` WHERE account = ? AND target = ?`, this.Account, this.Target)
	case KEY_Sight_Level:
		row = db.QueryRow(`SELECT id FROM `+this.TableName+` WHERE level = ?`, this.Level)
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

func (this *Sight) Load(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		return ErrRowNotFound
	}
	row := db.QueryRow(`SELECT * FROM `+this.TableName+` WHERE id = ?`, this.ID)
	return row.Scan(&this.ID, &this.Account, &this.Target, &this.Level, &this.UpdatedAt)
}

func (this *Sight) Save(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		_, err := db.Exec(`INSERT INTO `+this.TableName+` (level) VALUES (?)`, Sight_Unsaved)
		if err != nil {
			return err
		}
		level := this.Level
		this.Level = Sight_Unsaved
		this.ID = this.Find(KEY_Sight_Level)
		this.Level = level
	}
	// Update last updated time
	// stackoverflow.com/q/2218662
	_, err := db.Exec(`UPDATE `+this.TableName+` SET account = ?, target = ?, level = ?, updated_at = datetime('now') WHERE id = ?`, this.Account, this.Target, this.Level, this.ID)
	return err
}

func SightCount(tbl string, tgtid int) map[int]int {
	ct := map[int]int{}
	rows, err := db.Query(`SELECT level FROM `+tbl+` WHERE target = ?`, tgtid)
	if err != nil {
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		var level int
		if rows.Scan(&level) == nil {
			ct[level]++
		}
	}
	return ct
}
