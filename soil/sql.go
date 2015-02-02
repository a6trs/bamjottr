package soil

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"time"
)

type Storable interface {
	Find(key int) error
	Load(key int) error
	Save(key int) error
}

var ErrRowNotFound = errors.New("Storable.Load: Not found")
var ErrMembersNotFound = errors.New("GetMembers: No members found u^u")
var ErrAccountsNotFound = errors.New("FindAccounts: No accounts found -~-")

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
	err = init_Project()
	if err != nil {
		return err
	}
	err = init_ProjectTeamMembership()
	if err != nil {
		return err
	}
	err = init_Post()
	if err != nil {
		return err
	}
	err = init_Sight("projects")
	if err != nil {
		return err
	}
	err = init_Sight("posts")
	if err != nil {
		return err
	}
	err = init_Notification()
	if err != nil {
		return err
	}
	err = init_Invitation()
	if err != nil {
		return err
	}
	rand.Seed(time.Now().UnixNano())
	return nil
}
