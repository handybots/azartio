package handler

import (
	"time"

	"github.com/jackc/pgx/pgtype"

	tele "gopkg.in/tucnak/telebot.v3"
)

type Balance struct {
	User    *tele.User
	Perks   []pgtype.Varchar
	Balance int64
}

func (h handler) OnBalance(c tele.Context) error {
	defer c.Respond()

	user, err := h.db.Users.ByID(c.Sender())
	if err != nil {
		return err
	}

	balance := Balance{
		User: c.Sender(),
		Perks: user.Perks.Elements,
		Balance: user.Balance,
	}

	if c.Message().Private() {
		if c.Callback() != nil {
			_, err = h.b.Edit(
				c.Message(),
				h.lt.Text(c, "balance", balance),
				h.lt.Markup(c, "roulette"),
				)
			return err
		}
	}

	msg, err := h.b.Send(c.Chat(), h.lt.Text(c, "balance", balance))
	if err != nil {
		return err
	}

	time.Sleep(5 * time.Second)
	return h.b.Delete(msg)
}

func (h *handler) OnBonuses(c tele.Context) error {
	if c.Callback() != nil {
		return c.Edit(h.lt.Text(c, "bonuses", c.Sender().Recipient()), tele.NoPreview)
	}
	return c.Send(h.lt.Text(c, "bonuses", c.Sender().Recipient()), tele.NoPreview)
}
