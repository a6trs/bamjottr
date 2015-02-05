package soil

import (
	"database/sql"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"
)

type Account struct {
	ID       int
	Name     string
	Email    string
	Password []byte
	// The time when the notifications page was last visited
	LastRead time.Time
}

func init_Account() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(32),
		email VARCHAR(64),
		password VARCHAR(64),
		lastread DATETIME
	)`)
	return err
}

const (
	KEY_Account_ID = iota
	KEY_Account_Name
	KEY_Account_Email
	Account_PswdChangeMark = 233
)

func (this *Account) Find(key int) int {
	result := -1
	var row *sql.Row
	switch key {
	case KEY_Account_ID:
		row = db.QueryRow(`SELECT id FROM accounts WHERE id = ?`, this.ID)
	case KEY_Account_Name:
		row = db.QueryRow(`SELECT id FROM accounts WHERE name = ?`, this.Name)
	case KEY_Account_Email:
		row = db.QueryRow(`SELECT id FROM accounts WHERE email = ?`, this.Email)
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

func (this *Account) Load(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		return ErrRowNotFound
	}
	row := db.QueryRow(`SELECT * FROM accounts WHERE id = ?`, this.ID)
	return row.Scan(&this.ID, &this.Name, &this.Email, &this.Password, &this.LastRead)
}

func (this *Account) Save(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		passhash, err := bcrypt.GenerateFromPassword(this.Password, 10)
		if err != nil {
			return err
		}
		_, err = db.Exec(`INSERT INTO accounts (name, password) VALUES (?, ?)`, this.Name, passhash)
		if err != nil {
			return err
		}
		this.ID = this.Find(KEY_Account_Name)
	}
	changingPswd := this.Password[0] == Account_PswdChangeMark
	var err error
	if changingPswd {
		_, err = db.Exec(`UPDATE accounts SET name = ?, email = ?, lastread = ?, password = ? WHERE id = ?`, this.Name, this.Email, this.LastRead, this.Password[1:], this.ID)
	} else {
		_, err = db.Exec(`UPDATE accounts SET name = ?, email = ?, lastread = ? WHERE id = ?`, this.Name, this.Email, this.LastRead, this.ID)
	}
	return err
}

func (this *Account) MatchesPassword(pwd []byte) bool {
	err := bcrypt.CompareHashAndPassword(this.Password, pwd)
	return (err == nil)
}

func (this *Account) ChangePassword(pwd string) {
	passhash, err := bcrypt.GenerateFromPassword([]byte(pwd), 10)
	// stackoverflow.com/q/16248241
	if err == nil {
		this.Password = append([]byte{Account_PswdChangeMark}, passhash...)
	}
}

// Just a helper function, can be replaced by Find() and Save() calls.
// This simplifies the process of updating `LastRead` field a lot.
func UpdateLastReadTime(aid int) error {
	_, err := db.Exec(`UPDATE accounts SET lastread = datetime('now') WHERE id = ?`, aid)
	return err
}

// See above, faster than {{(account <AID>).Name}}
func AccountName(aid int) string {
	row := db.QueryRow(`SELECT name FROM accounts WHERE id = ? LIMIT 1`, aid)
	var s string
	if row.Scan(&s) == nil {
		return s
	} else {
		return "..."
	}
}

type InvitationState struct {
	Account *Account
	Invited bool
}

// TODO: Make these fmt.Sprintf calls look better
func FindAccounts(prjid, exclusion int, q string) ([]InvitationState, error) {
	ret := []InvitationState{}
	// Retrieve all accounts that have been invited
	rows1, err := db.Query(fmt.Sprintf(`SELECT accounts.* FROM accounts INNER JOIN invitations ON project = ? AND accounts.id = invitations.receiver AND name LIKE '%%%s%%'`, q), prjid)
	if err != nil {
		return nil, err
	}
	defer rows1.Close()
	list := []string{}
	for rows1.Next() {
		a := &Account{}
		if rows1.Scan(&a.ID, &a.Name, &a.Email, &a.Password, &a.LastRead) == nil {
			list = append(list, strconv.Itoa(a.ID))
			ret = append(ret, InvitationState{a, true})
		}
	}
	list = append(list, strconv.Itoa(exclusion))
	// Find all accounts
	rows2, err := db.Query(fmt.Sprintf(`SELECT * FROM accounts WHERE id NOT IN (%s) AND name LIKE '%%%s%%' LIMIT 10`, strings.Join(list, ","), q))
	if err != nil {
		return nil, err
	}
	defer rows2.Close()
	for rows2.Next() {
		a := &Account{}
		if rows2.Scan(&a.ID, &a.Name, &a.Email, &a.Password, &a.LastRead) == nil {
			ret = append(ret, InvitationState{a, false})
		}
	}
	return ret, nil
}
