package handler

import (
	tele "gopkg.in/tucnak/telebot.v3"
	"time"
)

func (h *handler) OnBonuses(c tele.Context) error {
	if !c.Message().Private() {
		msg, err := h.b.Reply(
			c.Message(),
			h.lt.Text(c, "bonuses", c.Sender().ID),
			h.lt.Markup(c, "bonuses"),
			tele.NoPreview,
		)

		time.AfterFunc(h.lt.Duration("delete_delay"), func() {
			c.Delete()
			h.b.Delete(msg)
		})

		return err
	}

	return c.Send(
		h.lt.Text(c, "bonuses", c.Sender().ID),
		h.lt.Markup(c, "bonuses"),
		tele.NoPreview,
	)
}

func (h handler) OnBonusDaily(c tele.Context) error {
	used, err := h.db.Users.IsLastBonusUsed(c.Sender())
	if err != nil {
		return err
	}

	if used {
		return c.Respond(&tele.CallbackResponse{
			Text: h.lt.Text(c, "bonus_used"),
		})
	}

	bonus := h.lt.Int64("bonuses.daily")
	if err = h.chargeBonus(c.Sender(), &bonus); err != nil {
		return err
	}
	if err := h.db.Users.UseBonus(c.Sender()); err != nil {
		return err
	}

	return c.Respond(&tele.CallbackResponse{
		Text: h.lt.Text(c, "bonus", bonus),
	})
}

func (h handler) OnBonusSponsor(c tele.Context) error {
	if h.db.Users.Subscribed(c.Sender()) {
		return c.Respond(&tele.CallbackResponse{
			Text:      h.lt.Text(c, "has_subscribed"),
			ShowAlert: true,
		})
	}

	if !h.subscribedOnSponsor(c.Sender()) {
		return c.Respond(&tele.CallbackResponse{
			Text:      h.lt.Text(c, "not_subscribed"),
			ShowAlert: true,
		})
	}

	bonus := h.lt.Int64("bonuses.sponsor")
	if err := h.chargeBonus(c.Sender(), &bonus); err != nil {
		return err
	}

	if err := h.db.Users.SetSubscribed(c.Sender(), true); err != nil {
		return err
	}

	return c.Send(h.lt.Text(c, "bonus", bonus))
}

func (h handler) subscribedOnSponsor(r tele.Recipient) bool {
	sponsor := h.lt.ChatID("sponsor_chat")

	member, err := h.b.ChatMemberOf(sponsor, r)
	if err != nil {
		return false
	}

	switch member.Role {
	case tele.Creator, tele.Administrator, tele.Member:
		return true
	default:
		return false
	}
}

func (h handler) chargeBonus(r tele.Recipient, bonus *int64) error {
	user, err := h.db.Users.ByID(r)
	if err != nil {
		return err
	}

	if h.dons.Scope("double_bonuses", user.Perks()...) {
		*bonus *= 2
	}

	return h.db.Users.Charge(*bonus, r)
}
