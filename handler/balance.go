package handler

import (
	"fmt"
	"time"

	tele "gopkg.in/tucnak/telebot.v3"
)

func (h *handler) OnBalance(c tele.Context) error {
	defer c.Respond()
	balance, err := h.db.Users.Balance(c.Sender())
	if err != nil {
		return err
	}
	text := fmt.Sprintf(`<a href="tg://user?id=%s">%s</a>, твой баланс: %d 💸`, c.Sender().Recipient(), c.Sender().FirstName, balance)
	if c.Message().Private() {
		if c.Callback() != nil {
			_, err = h.b.Edit(c.Message(), text, h.lt.Markup(c, "roulette"))
			return err
		}
	}
	m, err := h.b.Send(c.Chat(), text)
	time.Sleep(5 * time.Second)
	if err := h.b.Delete(m); err != nil {
		h.b.OnError(err, c)
	}
	return err
}
