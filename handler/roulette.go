package handler

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/handybots/azartio/azartio"
	"github.com/handybots/azartio/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

// States of group
const (
	rollingState = "rolling"
	idleState    = "none"
)

type key struct {
	UserID int
	Sign   string
}
type Bets struct {
	Bets  map[key]*storage.Bet
	Chats map[int]*tele.Chat
}

func (h *handler) OnBet(c tele.Context) error {
	defer c.Respond()
	group, err := h.db.Groups.ByID(c.Chat())
	if err != nil {
		return err
	}
	reservedMsg := tele.StoredMessage{
		MessageID: strconv.Itoa(int(group.MessageID)),
		ChatID:    c.Chat().ID,
	}
	var amount int64
	if c.Data() != "x2" {
		amount, err = strconv.ParseInt(c.Data(), 10, 64)
		if err != nil {
			return err
		}
	} else {
		bets, err := h.db.Bets.NotDoneByUserID(c.Sender())
		if err != nil {
			if err == sql.ErrNoRows {
				c.Respond(&tele.CallbackResponse{Text: h.lt.Text(c, "did_not_bet")})
			} else {
				return err
			}
		}
		betsx, err := h.collapseBets(bets)
		for _, bet := range betsx.Bets {
			if bet.Sign == c.Callback().Unique {
				log.Println(bet.Amount * 2)
				amount = bet.Amount * 2
				break
			}
		}
		if amount == 0 {
			c.Respond(&tele.CallbackResponse{Text: h.lt.Text(c, "did_not_bet")})
			return nil
		}
	}
	balance, err := h.db.Users.Balance(c.Sender())
	if err != nil {
		return err
	}

	if balance-amount <= 0 {
		err := c.Respond(&tele.CallbackResponse{
			Text: h.lt.Text(c, "not_enough_money"),
		})
		return err
	}

	if _, ok := azartio.Colors[c.Callback().Unique]; ok == false {
		return errors.New("handler: unknown sign in callback unique")
	}

	h.db.Users.Charge(-amount, c.Sender())

	bet := azartio.NewBet(c.Callback().Unique, amount, c.Sender().ID)

	err = h.db.Bets.Create(bet, c.Chat())
	if err != nil {
		c.Respond(&tele.CallbackResponse{
			Text: err.Error(),
		})
		return err
	}

	bets, err := h.db.Bets.NotDoneByChat(c.Chat())
	if err != nil {
		return err
	}

	betsx, err := h.collapseBets(bets)
	if err != nil {
		return err
	}

	h.b.Edit(
		reservedMsg,
		h.lt.Text(c, "bets", betsx),
		h.lt.Markup(c, "roulette"),
	)
	return nil
}

func (h *handler) OnRoulette(c tele.Context) error {
	h.b.Delete(c.Message())
	msg, err := h.b.Send(c.Chat(),
		h.lt.Text(c, "menu"),
		h.lt.Markup(c, "roulette"),
	)
	h.b.Pin(msg, tele.Silent)

	exists, err := h.db.Groups.Exists(c.Chat())
	if err != nil {
		return err
	}
	if exists {
		group, err := h.db.Groups.ByID(c.Chat())
		if err != nil {
			return err
		}
		h.b.Delete(
			tele.StoredMessage{
				MessageID: strconv.Itoa(int(group.MessageID)),
				ChatID:    c.Chat().ID},
		)

	} else {
		err := h.db.Groups.Create(c.Chat())
		if err != nil {
			return err
		}
	}
	h.db.Groups.UpdateMessage(c.Chat(), msg.ID)

	return err
}

func (h *handler) OnGo(c tele.Context) error {
	defer c.Respond()
	group, err := h.db.Groups.ByID(c.Chat())
	if err != nil {
		return err
	}
	reservedMsg := tele.StoredMessage{
		MessageID: strconv.Itoa(int(group.MessageID)),
		ChatID:    c.Chat().ID,
	}
	bets, err := h.db.Bets.NotDoneByChat(c.Chat())
	if err != nil {
		return err
	}

	j, err := h.collapseBets(bets)
	if err != nil {
		return err
	}

	azartioBets := make([]*azartio.Bet, 0, len(j.Bets))
	for _, bet := range j.Bets {
		azartioBets = append(azartioBets, bet.ToAzartioBet(bet.UserID))
	}
	type Result struct {
		Results map[key]*azartio.RollResult
		Chats   map[int]*tele.Chat
		WonSign string
	}
	results, err := h.c.RollMany(azartioBets)
	if err != nil {
		return err
	}

	jointResults := make(map[key]*azartio.RollResult)
	for _, result := range results {
		if _, ok := j.Chats[result.Bet.UserID]; !ok {
			chat, err := h.b.ChatByID(strconv.Itoa(result.Bet.UserID))
			if err != nil {
				return err
			}
			j.Chats[result.Bet.UserID] = chat
		}

		k := key{result.Bet.UserID, result.Bet.Sign}
		jointResults[k] = result
		if result.Won {
			err := h.db.Users.Charge(+result.Amount, j.Chats[result.Bet.UserID])
			if err != nil {
				return err
			}
		}
		err := h.db.Bets.MakeDone(result, j.Chats[result.Bet.UserID])
		if err != nil {
			return err
		}
	}

	err = h.db.Groups.UpdateState(c.Chat(), rollingState)
	if err != nil {
		return err
	}

	_, err = h.b.Edit(reservedMsg, h.lt.Text(c, "rolling"))
	if err != nil {
		return err
	}
	a, b := h.lt.Text(c, "roll_result", Result{
		WonSign: results[0].WonSign,
		Results: jointResults,
		Chats:   j.Chats,
	}), h.lt.Markup(c, "roulette")

	time.AfterFunc(3*time.Second, func() {
		_, err := h.b.Edit(reservedMsg, a, b)
		if err != nil {
			h.b.OnError(err, c)
		}
		err = h.db.Groups.UpdateState(c.Chat(), idleState)
		if err != nil {
			h.b.OnError(err, c)
		}
	},
	)

	return err
}

func (h *handler) OnPinned(c tele.Context) error {
	if c.Sender().ID == h.b.Me.ID {
		h.b.Delete(c.Message())
	}
	return nil
}

func (h *handler) collapseBets(bets []*storage.Bet) (*Bets, error) {
	chats := make(map[int]*tele.Chat)
	joint := make(map[key]*storage.Bet)

	for _, bet := range bets {
		if _, ok := chats[bet.UserID]; !ok {
			chat, err := h.b.ChatByID(strconv.Itoa(bet.UserID))
			if err != nil {
				return nil, err
			}
			chats[bet.UserID] = chat
		}

		k := key{bet.UserID, bet.Sign}
		if b, ok := joint[k]; ok {
			b.Amount += bet.Amount
		} else {
			joint[k] = bet
		}
	}
	return &Bets{Bets: joint, Chats: chats}, nil
}
