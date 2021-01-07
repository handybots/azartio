package handler

import (
	"math/rand"

	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) OnBonus(c tele.Context) error {
	used, err := h.db.Users.IsLastBonusUsed(c.Sender())
	if err != nil {
		return err
	}
	if used {
		_, err := h.b.Send(c.Chat(), h.lt.Text(c, "bonus_used"))
		return err
	}

	bonus := rand.Intn(3000)

	if err = h.db.Users.Charge(int64(bonus), c.Sender()); err != nil {
		return err
	}
	if err := h.db.Users.UseBonus(c.Sender()); err != nil {
		return err
	}

	_, err = h.b.Send(c.Chat(), h.lt.Text(c, "bonus", bonus))
	return err
}
