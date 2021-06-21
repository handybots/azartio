package handler

import (
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
