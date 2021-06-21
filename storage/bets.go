package storage

import (
	"github.com/handybots/azartio/azartio"
	"github.com/jmoiron/sqlx"
)

type (
	BetsStorage interface {
		ByID(id int) (bet Bet, _ error)
		Create(bet *azartio.Bet, chatID Chat) error
		ByUserID(chat Chat) (bets []Bet, _ error)
		NotDoneByUserID(user Chat, chat Chat) (bets []Bet, _ error)
		MakeDoneByChat(result *azartio.RollResult, user Chat, chat Chat) error
		NotDoneByChat(chat Chat) (bets []Bet, _ error)
	}

	Bets struct {
		*sqlx.DB
	}

	Bet struct {
		ID     int    `db:"id" sq:"id,omitempty"`
		UserID int    `sq:"user_id,omitempty"`
		ChatID int    `sq:"chat_id,omitempty"`
		Amount int64  `sq:"amount,omitempty"`
		Sign   string `sq:"sign,omitempty"`
		Won    bool   `sq:"won,omitempty"`
		Done   bool   `sq:"done,omitempty"`
	}
)

func (db *Bets) Create(bet *azartio.Bet, chat Chat) error {
	const q = `insert into bets (user_id, chat_id, amount, sign) values ($1,$2,$3,$4)`
	_, err := db.Exec(q, bet.UserID, chat.Recipient(), bet.Amount, bet.Sign)
	return err
}

func (db *Bets) ByID(id int) (bet Bet, _ error) {
	const q = `select 1 from bets where id = $1`
	return bet, db.Get(&bet, q, id)
}

func (db *Bets) ByUserID(chat Chat) (bets []Bet, _ error) {
	const q = `select * from bets where user_id = $1`
	return bets, db.Select(&bets, q, chat.Recipient())
}

func (db *Bets) NotDoneByUserID(user Chat, chat Chat) (bets []Bet, _ error) {
	const q = `select * from bets where user_id = $1 and chat_id = $2 and done = false`
	return bets, db.Select(&bets, q, user.Recipient(), chat.Recipient())
}

func (db *Bets) MakeDoneByChat(result *azartio.RollResult, user Chat, chat Chat) error {
	const q = `update bets set (won, amount, done) = ($2, $3, true) where user_id = $1 and won = false and sign = $4 and chat_id = $5 and done = false`
	_, err := db.Exec(q, user.Recipient(), result.Won, result.Amount, result.Bet.Sign, chat.Recipient())
	return err
}

func (db *Bets) NotDoneByChat(chat Chat) (bets []Bet, _ error) {
	const q = `select * from bets where chat_id = $1 and done = false`
	return bets, db.Select(&bets, q, chat.Recipient())
}
