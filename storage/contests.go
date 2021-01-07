package storage

import (
	"github.com/jackc/pgx/pgtype"
	"github.com/jmoiron/sqlx"
)

type (
	Contests struct {
		*sqlx.DB
	}

	Contest struct {
		CreatorID    int                 `db:"creator_id" sq:"creator_id,omitempty"`
		ChatID       int                 `db:"chat_id" sq:"chat_id,omitempty"`
		ID           int64               `db:"id" sq:"id, omitempty"`
		Amount       int64               `sq:"amount,omitempty"`
		Done         bool                `sq:"done,omitempty"`
		WinnerID     int64               `db:"winner_id" sq:"winner_id,omitempty"`
		Participants pgtype.VarcharArray `sq:"participants"`
		Canceled     bool                `sq:"canceled,omitempty"`
	}

	ContestsStorage interface {
		Create(creator Chat, chat Chat, amount int64) error
		NotDoneByUserID(userID string) (c Contest, _ error)
		AddParticipant(creator_id string, user Chat) error
		ByUserID(user Chat) (c []Contest, _ error)
		MarkDoneByUserID(user Chat) error
	}
)

func (db *Contests) Create(creator Chat, chat Chat, amount int64) error {
	const q = `insert into contests (creator_id, amount, chat_id, participants) values ($1,$2,$3, ARRAY['1'])`
	_, err := db.Exec(q, creator.Recipient(), amount, chat.Recipient())
	return err
}

func (db *Contests) NotDoneByUserID(userID string) (c Contest, _ error) {
	const q = `select * from contests where creator_id = $1 and not done`
	return c, db.Get(&c, q, userID)
}

func (db Contests) ExistsByUserID(userID string) (has bool, _ error) {
	const q = `select exists(select * from contests where creator_id = $1 and not done)`
	return has, db.Get(&has, q, userID)
}

func (db Contests) AddParticipant(creator_id string, user Chat) error {
	const q = `update contests set participants = array_append(participants, $2 ) where creator_id = $1 and not done`
	_, err := db.Exec(q, creator_id, user.Recipient())
	return err
}

func (db *Contests) ByUserID(user Chat) (c []Contest, _ error) {
	const q = `select * from contests where creator_id = $1`
	return c, db.Select(&c, q, user.Recipient())
}

func (db *Contests) MarkDoneByUserID(user Chat) error {
	const q = `update contests set done = true where creator_id = $1`
	_, err := db.Exec(q, user.Recipient())
	return err
}
