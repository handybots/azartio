package handler

import (
	tele "gopkg.in/tucnak/telebot.v3"
)

func (h *handler) OnPerks(c tele.Context) error {
	_, err := h.b.Edit(c.Message(), h.lt.Text(c, "perks"), h.lt.Markup(c, "perks_markup"))
	if err != nil {
		_, err = h.b.Send(c.Chat(), h.lt.Text(c, "perks"), h.lt.Markup(c, "perks_markup"))
	}
	return err

}

func (h *handler) OnPerk(c tele.Context) error {
	type BuyPerks struct {
		URL   string
		Price int
	}
	perk := h.lt.Text(c, c.Data()+"_desc")
	if perk == "" {
		h.b.Edit(c.Message(), "Not available now", h.lt.Markup(c, "back_to_perk"))
		return nil
	}
	d := h.d.Get(c.Data())
	i := d.Int("price")
	_, err := h.b.Edit(c.Message(), perk, h.lt.Markup(c, "buy_perk_menu", BuyPerks{"https://www.google.com/", i}))
	return err

}
