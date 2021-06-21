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

		time.AfterFunc(10*time.Second, func() {
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
	if err = h.db.Users.Charge(bonus, c.Sender()); err != nil {
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
	sponsor := h.lt.ChatID("sponsor_chat")
	member, err := h.b.ChatMemberOf(sponsor, c.Sender())
	if err != nil {
		return err
	}

	switch member.Role {
	case tele.Restricted, tele.Kicked, tele.Left:
		return c.Respond(&tele.CallbackResponse{
			Text:      h.lt.Text(c, "not_subscribed"),
			ShowAlert: true,
		})
	}

	bonus := h.lt.Int64("bonuses.sponsor")
	if err := h.db.Users.Charge(bonus, c.Sender()); err != nil {
		return err
	}

	return c.Send(h.lt.Text(c, "bonus", bonus))
}
