package handler

import (
	tele "gopkg.in/tucnak/telebot.v3"
	"strconv"
)

func (h handler) AdminBalance(c tele.Context) error {
	amount, _ := strconv.ParseInt(c.Message().Payload, 10, 64)
	return h.db.Users.Charge(amount, c.Recipient())
}

func (h handler) AdminPerk(c tele.Context) error {
	perk := c.Message().Payload // TODO: c.Data()

	don := h.dons.Get(perk)
	if don.Level == 0 {
		return nil
	}

	return h.db.Users.AddPerk(c.Recipient(), perk)
}
