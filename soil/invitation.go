package soil

// RSVP :) Should it be RSTP (s'il te - instead of vous - plai^t) here...

import (
	"math/rand"
	"strconv"
)

type Invitation struct {
	ID       int
	Project  int
	Receiver int
	Token    int64
}

func init_Invitation() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS invitations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project INTEGER,
		receiver INTEGER,
		token INTEGER UNIQUE,
		FOREIGN KEY(project) REFERENCES projects(id)
		FOREIGN KEY(receiver) REFERENCES accounts(id)
	)`)
	return err
}

func InvitationLink(project, receiver int) string {
	token := rand.Int63()
	_, err := db.Exec(`INSERT INTO invitations (project, receiver, token) VALUES (?, ?, ?)`, project, receiver, token)
	if err == nil {
		return "/answer_invitation/"+strconv.FormatInt(token, 36)
	} else {
		// fmt.Println("Oops.", err.Error())
		return ""
	}
}

func InvitationByToken(token int64) *Invitation {
	row := db.QueryRow(`SELECT * FROM invitations WHERE token = ?`, token)
	rsvp := &Invitation{}
	if row.Scan(&rsvp.ID, &rsvp.Project, &rsvp.Receiver, &rsvp.Token) != nil {
		return nil
	}
	return rsvp
}
