package handler

import (
	"errors"
	"github.com/handybots/azartio/azartio"
	tele "gopkg.in/tucnak/telebot.v3"
	"strconv"
	"strings"
)



func (h *handler) OnBet(c tele.Context) error {
	defer c.Respond()
	group, err := h.db.Groups.ByID(c.Chat())
	if err != nil{
		return err
	}
	reservedMsg := tele.StoredMessage{
		MessageID: strconv.Itoa(int(group.MessageID)),
		ChatID: c.Chat().ID}
	var amount int64 = 100
	var sign string
	sign = c.Data()
	if c.Callback() == nil{
		var s string
		m := strings.Split(c.Message().Text, " ")
		if len(m) == 1{
			s = m[0][:len(m[0])-2]
			sign = azartio.TranslateBet(m[0][len(m[0])-2:])
			}else{
				s = m[0]
				sign = azartio.TranslateBet(m[1])
				if sign == ""{
					sign = azartio.TranslateBet(m[2])
				}
			}
			var err error

			amount, err = strconv.ParseInt(s, 10, 64)
			if err != nil{
				return err
			}
			if sign == ""{
				return errors.New("debil")
			}
	}


	balance, err := h.db.Users.Balance(c.Sender())
	if err != nil{
		return err
	}
	if balance - amount <= 0 {
		err := c.Respond(&tele.CallbackResponse{Text:h.lt.Text(c, "not_enough_money")})
		return err
	}
	bet := azartio.NewBet(sign, amount, c.Sender())
	err = h.db.Bets.Create(bet, c.Chat())
	if err != nil{
		c.Respond(&tele.CallbackResponse{Text: err.Error()})
		return err
	}

	h.b.Edit(reservedMsg, h.lt.Text(c, "bets", ))
	return nil
}


func (h *handler) OnRoulette(c tele.Context) error {
	h.b.Delete(c.Message())
	msg, err := h.b.Send(c.Chat(),
		h.lt.Text(c, "menu"),
		h.lt.MarkupLocale("double", "roulette"),
	)
	h.b.Pin(msg, tele.Silent)

	if !c.Message().Private() {
		exists, err := h.db.Groups.Exists(c.Chat())
		if err != nil {
			return err
		}
		if exists {
			group, err := h.db.Groups.ByID(c.Chat())
			if err != nil{
				return err
			}
			h.b.Delete(
				tele.StoredMessage{
					MessageID: strconv.Itoa(int(group.MessageID)),
					ChatID: c.Chat().ID},
			)

		}else{
			err := h.db.Groups.Create(c.Chat())
			if err != nil {
				return err
			}
		}
		h.db.Groups.UpdateMessage(c.Chat(), msg.ID)

	}

	return err
}

func (h *handler) OnPinned(c tele.Context)error {
	if c.Sender().ID == h.b.Me.ID{
		h.b.Delete(c.Message())
	}
	return nil
}