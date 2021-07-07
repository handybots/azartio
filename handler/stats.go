package handler

import (
	"time"

	tele "gopkg.in/tucnak/telebot.v3"
)

type Stats struct {
	Balance int64
	Perks   int
	Bets    int
	Friends int
}

func (h handler) OnStats(c tele.Context) error {
	defer c.Respond()

	user, err := h.db.Users.ByID(c.Sender())
	if err != nil {
		return err
	}

	bets, err := h.db.Bets.Count(c.Sender())
	if err != nil {
		return err
	}

	friends, err := h.db.Users.Friends(c.Sender())
	if err != nil {
		return err
	}

	stats := Stats{
		Balance: user.Balance,
		Perks:   len(user.Perks()),
		Bets:    bets,
		Friends: friends,
	}

	if c.Callback() != nil {
		return c.Edit(
			h.lt.Text(c, "stats", stats),
			h.lt.Markup(c, "stats"),
		)
	}

	if !c.Message().Private() {
		msg, err := h.b.Reply(
			c.Message(),
			h.lt.Text(c, "stats", stats),
			h.lt.Markup(c, "stats"),
		)

		time.AfterFunc(h.lt.Duration("delete_delay"), func() {
			c.Delete()
			h.b.Delete(msg)
		})

		return err
	}

	return c.Send(
		h.lt.Text(c, "stats", stats),
		h.lt.Markup(c, "stats"),
	)
}

func (h handler) OnDeposit(c tele.Context) error {
	defer c.Respond()
	return c.Send(h.lt.Text(c, "coming_soon"), tele.NoPreview)
}
