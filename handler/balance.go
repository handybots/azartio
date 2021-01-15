package handler

import (
	"fmt"
	tele "gopkg.in/tucnak/telebot.v3"
)

func (h *handler) OnBalance(c tele.Context) error{
	balance, err := h.db.Users.Balance(c.Sender())
	if err != nil {
		return err
	}
	_, err = h.b.Send(c.Chat(), fmt.Sprintf(`<a href="tg://user?id=%s">%s</a>, твой баланс: %d 💸`, c.Sender().Recipient(), c.Sender().FirstName, balance))
	return err
}