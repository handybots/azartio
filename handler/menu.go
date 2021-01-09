package handler

import tele "gopkg.in/tucnak/telebot.v3"

func (h *handler) OnBet(c tele.Context) error {
	// TODO: bets bets bets
	defer c.Respond()
	return nil
}


func (h *handler) OnMenu(c tele.Context) error {
	_, err := h.b.Send(c.Chat(),
		h.lt.Text(c, "menu"),
		h.lt.Markup(c, "menu"),
	)
	return err
}