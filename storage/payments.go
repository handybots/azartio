package storage

import (
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/handybots/store"
	"github.com/jmoiron/sqlx"
)

type (
	PaymentsStorage interface {
		Create(payment Payment) (int, error)
		ByID(id int) (Payment, error)
		Pending(userID int, target, amount string) (Payment, error)
		Update(id int, profit string, payAt time.Time) error
	}

	Payments struct {
		*sqlx.DB
	}

	Payment struct {
		CreatedAt time.Time  `sq:"created_at,omitempty"`
		ID        int        `sq:"id,omitempty"`
		UserID    int        `sq:"user_id,omitempty"`
		Target    string     `sq:"target,omitempty"`
		Amount    string     `sq:"amount,omitempty"`
		Profit    string     `sq:"profit,omitempty"`
		PayAt     *time.Time `sq:"pay_at,omitempty"`
	}
)

func (p Payment) Payment() store.Payment {
	return store.Payment{
		ID:       p.ID,
		UserID:   p.UserID,
		Amount:   p.Amount,
		Currency: store.RUB,
		Target:   p.Target,
		Profit:   p.Profit,
		PayAt:    p.PayAt,
	}
}

func (p Payment) Payed() bool {
	return p.PayAt != nil
}

func (db *Payments) Create(p Payment) (id int, _ error) {
	q, args, err := sq.
		Insert("payments").
		SetMap(structs.Map(p)).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING id").
		ToSql()
	if err != nil {
		return 0, err
	}

	if err := db.QueryRow(q, args...).Scan(&id); err != nil {
		return 0, err
	}

	return
}

func (db *Payments) ByID(id int) (p Payment, _ error) {
	const q = `SELECT * FROM payments WHERE id=$1`
	return p, db.Get(&p, q, id)
}

func (db *Payments) Pending(userID int, target, amount string) (p Payment, _ error) {
	const q = `
		SELECT * FROM payments
		WHERE user_id=$1 AND target=$2
		AND amount=$3 AND pay_at IS NULL`

	return p, db.Get(&p, q, userID, target, amount)
}

func (db *Payments) Update(id int, profit string, payAt time.Time) error {
	const q = `UPDATE payments SET profit=$1, pay_at=$2 WHERE id=$3`
	_, err := db.Exec(q, profit, payAt, id)
	return err
}
