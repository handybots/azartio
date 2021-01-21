package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type (
	UsersStorage interface {
		Create(chat Chat, ref string) error
		Exists(chat Chat) (bool, error)
		Lang(chat Chat) (string, error)
		SetLang(chat Chat, lang string) error
		Charge(amount int64, chat Chat) error
		ByID(chat Chat) (usr User, _ error)
		Balance(chat Chat) (a int64, _ error)
		IsLastBonusUsed(chat Chat) (bool, error)
		UseBonus(chat Chat) error
		Leaderboard() (users []User, _ error)
	}

	Users struct {
		*sqlx.DB
	}

	User struct {
		CreatedAt time.Time `db:"created_at" sq:"created_at,omitempty"`
		UpdatedAt time.Time `sq:"updated_at,omitempty"`
		Balance   int64     `sq:"balance,omitempty"`
		ID        int       `db:"id" sq:"id,omitempty"`
		Lang      string    `sq:"lang,omitempty"`
		Ref       string    `sq:"ref"`
		LastBonus time.Time `sq:"last_bonus"`
	}

	Chat interface {
		Recipient() string
	}
)

func (db *Users) Create(chat Chat, ref string) error {
	const q = `INSERT INTO users (id, lang, ref) VALUES ($1, 'ru', $2)`
	_, err := db.Exec(q, chat.Recipient(), ref)
	return err
}

func (db *Users) Exists(chat Chat) (has bool, _ error) {
	const q = `SELECT EXISTS(SELECT 1 FROM users WHERE id=$1)`
	return has, db.Get(&has, q, chat.Recipient())
}

func (db *Users) Lang(chat Chat) (lang string, _ error) {
	const q = `SELECT lang FROM users WHERE id=$1`
	return lang, db.Get(&lang, q, chat.Recipient())
}

func (db *Users) SetLang(chat Chat, lang string) error {
	const q = `UPDATE users SET lang=$1 WHERE id=$2`
	_, err := db.Exec(q, lang, chat.Recipient())
	return err
}

func (db *Users) Charge(amount int64, chat Chat) error {
	const q = `update users set balance = balance + $1 where id = $2`
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.Exec(q, amount, chat.Recipient())
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (db *Users) ByID(chat Chat) (usr User, _ error) {
	const q = `select * from users where id = $1`
	return usr, db.Get(&usr, q, chat.Recipient())
}

func (db *Users) Balance(chat Chat) (a int64, _ error) {
	const q = `select balance from users where id = $1`
	return a, db.Get(&a, q, chat.Recipient())
}

func (db *Users) IsLastBonusUsed(chat Chat) (bool, error) {
	const q = `SELECT DATE_PART(
        'day',
        now()::timestamp - (
            select last_bonus
            from users
            where id = $1
        )::timestamp
    );`
	var dayDifference int
	err := db.Get(&dayDifference, q, chat.Recipient())
	if err != nil {
		return false, err
	}
	if dayDifference > 0 {
		return false, nil
	}
	return true, nil
}

func (db *Users) UseBonus(chat Chat) error {
	const q = `update users set last_bonus = now() where id = $1`
	_, err := db.Exec(q, chat.Recipient())
	return err
}

func (db *Users) Leaderboard() (users []User, _ error) {
	const q = `select * from users order by balance desc`
	return users, db.Select(&users, q)
}
