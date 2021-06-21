package handler

import (
	"github.com/sirupsen/logrus"
	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) LoggerFields(c tele.Context) logrus.Fields {
	f := make(logrus.Fields)

	switch {
	// Check callback first to avoid fetching its actual message.
	case c.Callback() != nil:
		f["event"] = "callback"
	case c.Message() != nil:
		f["event"] = "message"
	}

	user := c.Sender()
	if user != nil {
		f["user_id"] = user.Recipient()

		chat := c.Chat()
		if chat != nil && chat.Recipient() != user.Recipient() {
			f["chat_id"] = chat.Recipient()
		}
	}

	return f
}
