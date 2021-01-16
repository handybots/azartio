package handler

import (
	"errors"
	"fmt"
	"github.com/handybots/azartio/azartio"
	tele "gopkg.in/tucnak/telebot.v3"
	"strconv"
	"strings"
)



func (h *handler) OnBet(c tele.Context) error {
	defer c.Respond()
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
		_, err := h.b.Send(c.Chat(), h.lt.Text(c,"not_enough_money"))
		return err
	}
	bet := azartio.NewBet(sign, amount, c.Sender().ID)
	err = h.db.Bets.Create(bet, c.Chat())
	if err != nil{
		c.Respond(&tele.CallbackResponse{Text: err.Error()})
		return err
	}

	h.b.Send(c.Chat(),
		fmt.Sprintf(
			`<a href="tg://user?id=%s">%s</a>, ты поставил %s`,
			c.Sender().Recipient(),
			c.Sender().FirstName,
			bet.String()),
	)
	return nil
}


func (h *handler) OnMenu(c tele.Context) error {
	_, err := h.b.Send(c.Chat(),
		h.lt.Text(c, "menu"),
		h.lt.MarkupLocale("double", "menu"),
	)

	return err
}