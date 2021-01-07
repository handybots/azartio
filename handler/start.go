package handler

import (
	"errors"
	"log"

	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) OnStart(c tele.Context) error {
	if !c.Message().Private() {
		return errors.New("start from group")
	}

	var (
		chat = c.Sender()
		ref  = c.Message().Payload
	)

	exists, err := h.db.Users.Exists(chat)
	if err != nil {
		return err
	}

	if !exists {
		log.Println("Start from", chat.Recipient())
		if err := h.db.Users.Create(chat, ref); err != nil {
			return err
		}
		refBy, err := h.b.ChatByID(ref)
		if err == nil {
			err := h.db.Users.Charge(500, chat)
			if err == nil {
				h.b.Send(chat, h.lt.Text(c, "ref"))
			}

			h.db.Users.Charge(500, refBy)
			h.b.Send(refBy, h.lt.Text(c, "join_ref", chat.Recipient()))
		}
	}

	return c.Send(
		h.lt.Text(c, "start"),
		h.lt.Markup(c, "private_menu"),
	)
}

func (h handler) OnLang(c tele.Context) error {
	defer c.Respond()
	lang := c.Data()

	if locale, _ := h.lt.Locale(c); locale == lang {
		return nil
	}

	if err := h.db.Users.SetLang(c.Sender(), lang); err != nil {
		return err
	} else {
		h.lt.SetLocale(c, lang)
	}

	return c.Edit(
		h.lt.Text(c, "start"),
		h.lt.Markup(c, "lang"),
	)
}
