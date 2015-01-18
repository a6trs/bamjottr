package soil

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID       int
	Name     string
	Email    string
	Password []byte
}

func init_Account() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS accounts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(32),
		email VARCHAR(64),
		password VARCHAR(64)
	)`)
	return err
}

const (
	KEY_Account_ID = iota
	KEY_Account_Name
	KEY_Account_Email
)

func (this *Account) Find(key int) int {
	result := -1
	var stmt string
	switch key {
	case KEY_Account_ID:
		stmt = fmt.Sprintf(`SELECT id FROM accounts WHERE id = %d`, this.ID)
	case KEY_Account_Name:
		stmt = fmt.Sprintf(`SELECT id FROM accounts WHERE name = '%s'`, this.Name)
	case KEY_Account_Email:
		stmt = fmt.Sprintf(`SELECT id FROM accounts WHERE email = '%s'`, this.Email)
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

func (this *Account) Load(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		return ErrRowNotFound
	}
	row := db.QueryRow(fmt.Sprintf(`SELECT * FROM accounts WHERE id = %d`, this.ID))
	return row.Scan(&this.ID, &this.Name, &this.Email, &this.Password)
}

func (this *Account) Save(key int) error {
	this.ID = this.Find(key)
	if this.ID == -1 {
		passhash, err := bcrypt.GenerateFromPassword(this.Password, 10)
		if err != nil {
			return err
		}
		_, err = db.Exec(fmt.Sprintf(`INSERT INTO accounts (name, password) VALUES ('%s', %q)`, this.Name, passhash))
		if err != nil {
			return err
		}
		this.ID = this.Find(KEY_Account_Name)
	}
	_, err := db.Exec(fmt.Sprintf(`UPDATE accounts SET name = '%s', email = '%s' WHERE id = %d`, this.Name, this.Email, this.ID))
	return err
}

func (this *Account) MatchesPassword(pwd []byte) bool {
	err := bcrypt.CompareHashAndPassword(this.Password, pwd)
	return (err == nil)
}
