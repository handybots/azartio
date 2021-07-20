package handler

import (
	"log"
	"strings"

	"github.com/handybots/store"
	tele "gopkg.in/tucnak/telebot.v3"
)

type Payment struct {
	AnypayURL string
	EnotioURL string
	Amount    string
	Perk      string
}

func (h handler) OnPayment(p store.Payment) {
	logerr := func(v interface{}) {
		log.Println("handler/payment:", v)
	}

	payment, err := h.db.Payments.ByID(p.ID)
	if err != nil {
		logerr(err)
		return
	}

	if p.UserID != payment.UserID {
		logerr("bad request")
		return
	}

	if err := h.db.Payments.Update(p.ID, p.Profit, *p.PayAt); err != nil {
		logerr(err)
		return
	}

	switch {
	case strings.HasPrefix(p.Target, "perk:"):
		perk := strings.Split(p.Target, ":")[1]
		if h.dons.Get(perk).Level == 0 {
			logerr("perk " + perk + " does not exist")
			return
		}

		chat := tele.ChatID(p.UserID)
		if err := h.db.Users.AddPerk(chat, perk); err != nil {
			logerr(err)
			return
		}

		_, err = h.b.Send(chat, h.lt.TextLocale("ru", "payed_perk", perk))
		if err != nil {
			logerr(err)
			return
		}

		bonus := h.lt.Int64("bonuses.donate")
		if err := h.chargeBonus(chat, &bonus); err != nil {
			logerr(err)
		}
	}
}
