package handler

import (
	"errors"
	"fmt"
	"strings"

	"github.com/handybots/pkg/store"
	tele "gopkg.in/tucnak/telebot.v3"
)

type Payment struct {
	AnypayURL string
	EnotioURL string
	Amount    string
	Perk      string
}

func (h handler) OnPayment(p store.Payment) error {
	payment, err := h.db.Payments.ByID(p.ID)
	if err != nil {
		return err
	}

	if p.UserID != payment.UserID {
		return errors.New("handler/payment: bad request")
	}

	if err := h.db.Payments.Update(p.ID, p.Profit, *p.PayAt); err != nil {
		return err
	}

	switch {
	case strings.HasPrefix(p.Target, "perk:"):
		perk := strings.Split(p.Target, ":")[1]
		if h.dons.Get(perk).Level == 0 {
			return fmt.Errorf("handler/payment: perk %s does not exist", perk)
		}

		chat := tele.ChatID(p.UserID)
		if err := h.db.Users.AddPerk(chat, perk); err != nil {
			return err
		}

		_, err = h.b.Send(chat, h.lt.TextLocale("ru", "payed_perk", perk))
		return err
	}

	return nil
}
