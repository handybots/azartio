package handler

import (
	"encoding/json"
	"github.com/handybots/azartio/azartio"
	tele "gopkg.in/tucnak/telebot.v3"
)

func (h *handler) OnBet(c tele.Context) error {
	// TODO: bets bets bets
	bet := azartio.NewBet(c.Data(), 5, c.Sender().ID)
	err := h.db.Bets.Create(bet, c.Chat())
	if err != nil{
		c.Respond(&tele.CallbackResponse{Text: err.Error()})
		return err
	}
	result, err := h.c.Roll(bet)
	if err != nil{
		c.Respond(&tele.CallbackResponse{Text: err.Error()})
		return err
	}
	j, _ := json.Marshal(result)
	h.b.Send(c.Chat(), string(j))
	h.db.Bets.MakeDone(result, c.Sender())
	defer c.Respond()
	return nil
}


func (h *handler) OnMenu(c tele.Context) error {
	_, err := h.b.Send(c.Chat(),
		h.lt.Text(c, "menu"),
		h.lt.MarkupLocale("double", "menu"),
	)

	return err
}