package storage

import (
	"time"

	"github.com/jackc/pgx/pgtype"
	"github.com/jmoiron/sqlx"
)

type (
	UsersStorage interface {
		Create(chat Chat, ref string) error
		Exists(chat Chat) (bool, error)
		Charge(amount int64, chat Chat) error
		ByID(chat Chat) (usr User, _ error)
		Balance(chat Chat) (a int64, _ error)
		Friends(chat Chat) (int, error)
		IsLastBonusUsed(chat Chat) (bool, error)
		UseBonus(chat Chat) error
		Subscribed(chat Chat) bool
		SetSubscribed(chat Chat, s bool) error
		Leaderboard() (users []User, _ error)
		AddPerk(chat Chat, perk string) error
	}

	Users struct {
		*sqlx.DB
	}

	User struct {
		CreatedAt  time.Time           `sq:"created_at,omitempty"`
		UpdatedAt  time.Time           `sq:"updated_at,omitempty"`
		ID         int                 `sq:"id,omitempty"`
		Lang       string              `sq:"lang,omitempty"` // TODO: remove
		Ref        string              `sq:"ref,omitempty"`
		Balance    int64               `sq:"balance,omitempty"`
		LastBonus  time.Time           `sq:"last_bonus,omitempty"`
		Subscribed bool                `sq:"subscribed,omitempty"`
		PerksArray pgtype.VarcharArray `db:"perks" sq:"perks,omitempty"`
	}

	Chat interface {
		Recipient() string
	}
)

func (u User) HasPerk(perk string) bool {
	for _, p := range u.Perks() {
		if p == perk {
			return true
		}
	}
	return false
}

func (u User) Perks() (ps []string) {
	u.PerksArray.AssignTo(&ps)
	return
}

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
	_, err := db.Exec(q, amount, chat.Recipient())
	return err
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

func (db *Users) Subscribed(chat Chat) (sub bool) {
	const q = `select subscribed from users where id = $1`
	return db.Get(&sub, q, chat.Recipient()) == nil && sub
}

func (db *Users) SetSubscribed(chat Chat, sub bool) error {
	const q = `update users set subscribed = $2 where id = $1`
	_, err := db.Exec(q, chat.Recipient(), sub)
	return err
}

func (db *Users) Leaderboard() (users []User, _ error) {
	const q = `select * from users order by balance desc limit 10`
	return users, db.Select(&users, q)
}

func (db *Users) AddPerk(chat Chat, perk string) error {
	const q = `update users set perks = array_append(perks, $1) where id = $2`
	_, err := db.Exec(q, perk, chat.Recipient())
	return err
}

func (db *Users) Friends(chat Chat) (i int, _ error) {
	const q = `select count(*) from users where ref = $1`
	return i, db.Get(&i, q, chat.Recipient())
}
