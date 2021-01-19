package handler

import (
	"math/rand"

	tele "gopkg.in/tucnak/telebot.v3"
)

func (h *handler) OnBonus(c tele.Context) error {
	used, err := h.db.Users.IsLastBonusUsed(c.Sender())
	if err != nil {
		return err
	}
	if used {
		h.b.Send(c.Chat(), h.lt.Text(c, "bonus_used"))
		return nil
	}
	a := rand.Intn(3000)
	err = h.db.Users.Charge(int64(a), c.Sender())
	if err != nil {
		return err
	}
	err = h.db.Users.UseBonus(c.Sender())
	h.b.Send(c.Chat(), h.lt.Text(c, "bonus", a))

	return nil
}
