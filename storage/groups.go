package storage

import (
	"github.com/jmoiron/sqlx"
)

const (
	GroupStateNone    = "none"
	GroupStateRolling = "rolling"
)

type (
	GroupsStorage interface {
		Create(chat Chat) error
		UpdateMessage(chat Chat, MessageID int) error
		UpdateState(chat Chat, state string) error
		ByID(chat Chat) (group Group, _ error)
		Exists(chat Chat) (exists bool, _ error)
	}

	Groups struct {
		*sqlx.DB
	}

	Group struct {
		ID        int64  `db:"id" sq:"id,omitempty"`
		State     string `sq:"state,omitempty"`
		MessageID int64  `db:"message_id" sq:"message_id,omitempty"`
	}
)

func (db *Groups) Create(chat Chat) error {
	const q = `insert into groups (id) values($1)`
	_, err := db.Exec(q, chat.Recipient())
	return err
}

func (db *Groups) Exists(chat Chat) (has bool, _ error) {
	const q = `select exists(select 1 from groups where id = $1)`
	return has, db.Get(&has, q, chat.Recipient())
}

func (db *Groups) UpdateMessage(chat Chat, MessageID int) error {
	const q = `update groups set message_id = $1 where id = $2`
	_, err := db.Exec(q, MessageID, chat.Recipient())
	return err
}

func (db *Groups) UpdateState(chat Chat, state string) error {
	const q = `update groups set state = $1 where id = $2`
	_, err := db.Exec(q, state, chat.Recipient())
	return err
}

func (db *Groups) ByID(chat Chat) (group Group, _ error) {
	const q = `select * from groups where id = $1`
	return group, db.Get(&group, q, chat.Recipient())
}
