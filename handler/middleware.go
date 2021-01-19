package handler

import (
	"errors"

	tele "gopkg.in/tucnak/telebot.v3"
)

func (h *handler) Validate() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			exists, err := h.db.Users.Exists(c.Sender())
			if err != nil {
				return err
			}
			if !exists {
				h.b.Send(c.Chat(), h.lt.Text(c, "unregistered"))
				c.Respond()
				return errors.New("unregistered")
			}

			return next(c)
		}
	}
}
