package soil

import (
	"errors"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Storable interface {
	Find(key int) error
	Load(key int) error
	Save(key int) error
}

var ErrRowNotFound = errors.New("Storable.Load: Not found")

var db *sql.DB

func InitDatabase() error {
	var err error
	db, err = sql.Open("sqlite3", "data.db")
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	err = init_Account()
	if err != nil {
		return err
	}
	return nil
}
