package storage

import (
	"github.com/handybots/azartio/azartio"
	"github.com/jmoiron/sqlx"
)

type (
	BetsStorage interface {
		ByID(id int) (bet Bet, _ error)
		Create(bet azartio.Bet, chatID Chat) error
		ByUserID(chat Chat) (bets []Bet, _ error)
		NotDoneByUserID(chat Chat)(bets []Bet, _ error)
		MakeDone(result *azartio.RollResult, chat Chat) error
	}

	Bets struct {
		*sqlx.DB
	}

	Bet struct {
		ID int `db:"id" sq:"id,omitempty"`
		UserID int	`sq:"user_id,omitempty"`
		ChatID int	`sq:"chat_id,omitempty"`
		Amount int64 `sq:"amount,omitempty"`
		Sign string `sq:"sign,omitempty"`
		Won bool `sq:"won,omitempty"`
		Done bool `sq:"done,omitempty"`
	}
)

func (db *Bets) Create(bet azartio.Bet, chatID Chat) error {
	const q = `insert into bets (user_id, chat_id, amount, sign) values ($1,$2,$3,$4)`
	_, err := db.Exec(q,bet.ID, chatID.Recipient(), bet.Amount, bet.Sign)
	return err
}

func (db *Bets) ByID(id int) (bet Bet, _ error){
	const q = `select 1 from bets where id = $1`
	return bet, db.Get(&bet, q, id)
}

func (db *Bets) ByUserID(chat Chat) (bets []Bet, _ error){
	const q = `select * from bets where user_id = $1`
	return bets, db.Select(&bets,q,chat.Recipient())
}

func (db *Bets) NotDoneByUserID(chat Chat)(bets []Bet, _ error){
	const q = `select * from bets where user_id = $1 and done = false`
	return bets, db.Select(&bets,q,chat.Recipient())
}

func (db *Bets) MakeDone(result *azartio.RollResult, chat Chat) error{
	const q = `update bets set (won, amount, done) = ($2, $3, true) where user_id = $1 and won = false and sign = $4 and done = false`
	_, err := db.Exec(q, chat.Recipient(), result.Won, result.Amount, result.Bet.Sign)
	return err
}

// not done yet