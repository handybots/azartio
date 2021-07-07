package handler

import (
	"database/sql"
	"time"

	"github.com/handybots/azartio/storage"
	"github.com/handybots/store/anypay"
	"github.com/handybots/store/enotio"

	tele "gopkg.in/tucnak/telebot.v3"
)

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

		time.AfterFunc(h.lt.Duration("delete_delay"), func() {
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
	defer c.Respond()

	don := h.dons.Get(c.Data())
	if don.Level == 0 {
		return c.Edit(
			h.lt.Text(c, "coming_soon"),
			h.lt.Markup(c, "perks_back"),
			tele.NoPreview,
		)
	}

	user, err := h.db.Users.ByID(c.Sender())
	if err != nil {
		return err
	}

	if user.HasPerk(don.Name) {
		return c.Edit(
			h.lt.Text(c, "perk", don.Name),
			h.lt.Markup(c, "perks_back"),
		)
	}

	var (
		userID = c.Sender().ID
		target = "perk:" + don.Name
		amount = don.String("price")
	)

	p, err := h.db.Payments.Pending(userID, target, amount)
	if err == sql.ErrNoRows {
		p = storage.Payment{
			UserID: userID,
			Target: target,
			Amount: amount,
		}
		p.ID, err = h.db.Payments.Create(p)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	payment := Payment{
		Amount:    amount,
		Perk:      don.Name,
		AnypayURL: anypay.URL(p.Payment()),
		EnotioURL: enotio.URL(p.Payment()),
	}

	return c.Edit(
		h.lt.Text(c, "perk", don.Name),
		h.lt.Markup(c, "perk_buy", payment),
	)
}
