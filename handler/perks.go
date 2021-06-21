package handler

import (
	tele "gopkg.in/tucnak/telebot.v3"
	"time"
)

type BuyPerks struct {
	URL   string
	Price int
}

func (h *handler) OnPerks(c tele.Context) error {
	user, err := h.db.Users.ByID(c.Sender())
	if err != nil {
		return err
	}

	if c.Callback() != nil {
		return c.Edit(
			h.lt.Text(c, "perks", user.Perks()),
			h.lt.Markup(c, "perks"),
		)
	}

	if !c.Message().Private() {
		msg, err := h.b.Reply(
			c.Message(),
			h.lt.Text(c, "perks", user.Perks()),
			h.lt.Markup(c, "perks"),
		)

		time.AfterFunc(10*time.Second, func() {
			c.Delete()
			h.b.Delete(msg)
		})

		return err
	}

	return c.Send(
		h.lt.Text(c, "perks", user.Perks()),
		h.lt.Markup(c, "perks"),
	)
}

func (h *handler) OnPerk(c tele.Context) error {
	don := h.dons.Get(c.Data())
	if don.Level == 0 {
		return c.Edit(
			h.lt.Text(c, "coming_soon"),
			h.lt.Markup(c, "perks_back"),
			tele.NoPreview,
		)
	}

	return c.Edit(
		h.lt.Text(c, "perk", don.Name),
		h.lt.Markup(c, "perk_buy", don.Int("price")),
	)
}
