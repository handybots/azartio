package handler

import (
	"time"

	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) Validate(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		exists, err := h.db.Users.Exists(c.Sender())
		if err != nil {
			return err
		}

		if !exists && c.Callback() != nil {
			return c.Respond(&tele.CallbackResponse{
				Text: h.lt.Text(c, "unregistered"),
			})
		}

		return next(c)
	}
}

func (h handler) ApplyBonuses(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		chat := c.Sender()

		user, err := h.db.Users.ByID(chat)
		if err != nil {
			return err
		}

		if h.dons.Scope("auto_bonuses", user.Perks()...) {
			diff := int64(time.Now().Sub(user.LastBonus).Hours() / 24)
			if diff > 0 {
				bonus := diff * h.lt.Int64("bonuses.daily")
				if err := h.chargeBonus(chat, &bonus); err != nil {
					return err
				}
				if err := h.db.Users.UseBonus(chat); err != nil {
					return err
				}

				if err := c.Send(h.lt.Text(c, "bonus_hacker", bonus)); err != nil {
					return err
				}
			}

			if !h.subscribedOnSponsor(chat) && !h.db.Users.Subscribed(chat) {
				bonus := h.lt.Int64("bonuses.sponsor")
				if err := h.chargeBonus(chat, &bonus); err != nil {
					return err
				}
				if err := h.db.Users.SetSubscribed(chat, true); err != nil {
					return err
				}

				if err := c.Send(h.lt.Text(c, "bonus_hacker", bonus)); err != nil {
					return err
				}
			}
		}

		return next(c)
	}
}
