package storage

import "github.com/jmoiron/sqlx"

type (
	BetsStorage interface {
		ByID(id int) (bet Bet, _ error)
		Create(bet Bet) error
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

func (db *Bets) Create(userID Chat, chatID int,amount int64, sign string ) error {
	const q = `insert into bets (user_id, chat_id, amount, sign) values ($1,$2,$3,$4)`
	_, err := db.Exec(q,userID, chatID, amount, sign)
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

// not done yet