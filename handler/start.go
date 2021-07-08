package handler

import (
	"log"

	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) OnStart(c tele.Context) error {
	if !c.Message().Private() {
		return nil
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
		err := h.registerUser(ref, chat)
		if err != nil {
			return err
		}

	}

	return c.Send(
		h.lt.Text(c, "start"),
		h.lt.Markup(c, "private_menu"),
		tele.NoPreview,
	)
}

func (h *handler) registerUser(ref string, chat *tele.User) error {
	log.Println("Start from", chat.Recipient())
	if err := h.db.Users.Create(chat, ref); err != nil {
		return err
	}

	var (
		startBonus  = h.lt.Int64("bonuses.start")
		friendBonus = h.lt.Int64("bonuses.friend")
	)

	if err := h.db.Users.Charge(startBonus, chat); err != nil {
		return err
	}

	refBy, err := h.b.ChatByID(ref)
	if err == nil {
		if err := h.chargeBonus(chat, &friendBonus); err == nil {
			h.b.Send(chat, h.lt.Text(c, "ref", friendBonus))
		}
		if err := h.chargeBonus(refBy, &friendBonus); err == nil {
			defer h.b.Send(refBy, h.lt.TextLocale("ru", "join_ref", friendBonus))
		}
	}

}
