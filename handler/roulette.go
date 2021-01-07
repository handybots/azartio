package handler

import (
	"database/sql"
	"errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/pgtype"

	"github.com/handybots/azartio/azartio"
	"github.com/handybots/azartio/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

// States of group
const (
	stateIdle    = "none"
	stateRolling = "rolling"
)

type key struct {
	UserID int
	Sign   string
}

type Bets struct {
	Bets  map[key]*storage.Bet
	Chats map[int]*tele.Chat
}

func (h handler) OnBet(c tele.Context) error {
	defer c.Respond()

	group, err := h.db.Groups.ByID(c.Chat())
	if err != nil {
		return err
	}

	reservedMsg := tele.StoredMessage{
		MessageID: strconv.Itoa(int(group.MessageID)),
		ChatID:    c.Chat().ID,
	}

	var (
		amount int64
		sign   string
	)

	if c.Callback() == nil {
		var s string
		m := strings.Split(c.Message().Text, " ")

		if len(m) == 1 {
			s = m[0][:len(m[0])-2]
			sign = azartio.TranslateBet(m[0][len(m[0])-2:])
		} else {
			s = m[0]
			sign = azartio.TranslateBet(m[1])
			if sign == "" && len(m) > 1 {
				sign = azartio.TranslateBet(m[2])
			}
		}

		amount, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		if sign == "" {
			return nil
		}

		h.b.Delete(c.Message())
	} else {
		sign = c.Callback().Unique

		if c.Data() != "x2" {
			amount, err = strconv.ParseInt(c.Data(), 10, 64)
			if err != nil {
				return err
			}
		} else {
			bets, err := h.db.Bets.NotDoneByUserID(c.Sender(), c.Chat())
			if err != nil {
				if err == sql.ErrNoRows {
					c.Respond(&tele.CallbackResponse{Text: h.lt.Text(c, "did_not_bet")})
				} else {
					return err
				}
			}
			betsx, err := h.collapseBets(bets)
			for _, bet := range betsx.Bets {
				if bet.Sign == sign {
					amount = bet.Amount
					break
				}
			}
			if amount == 0 {
				c.Respond(&tele.CallbackResponse{Text: h.lt.Text(c, "did_not_bet")})
				return nil
			}
		}

	}

	balance, err := h.db.Users.Balance(c.Sender())
	if err != nil {
		return err
	}

	if balance-amount <= 0 {
		return c.Respond(&tele.CallbackResponse{
			Text: h.lt.Text(c, "not_enough_money"),
		})
	}

	if _, ok := azartio.Colors[sign]; ok == false {
		return errors.New("handler: unknown sign in callback unique")
	}

	if err := h.db.Users.Charge(-amount, c.Sender()); err != nil {
		return err
	}

	bet := azartio.NewBet(sign, amount, c.Sender().ID)
	if err = h.db.Bets.Create(bet, c.Chat()); err != nil {
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

	_, err = h.b.Edit(
		reservedMsg,
		h.lt.Text(c, "bets", betsx),
		h.lt.Markup(c, "roulette"),
	)
	return err
}

func (h handler) OnRoulette(c tele.Context) error {
	h.b.Delete(c.Message())
	if c.Message().Private() {
		usr, err := h.db.Users.ByID(c.Chat())
		if err != nil {
			return err
		}
		if time.Now().Sub(usr.CreatedAt) <= 24 {
			msg, err := h.b.Send(c.Chat(), h.lt.Text(c, "rules"))
			if err != nil {
				return err
			}
			time.Sleep(3 * time.Second)
			h.b.Delete(msg)
		}
	}
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

func (h handler) OnGo(c tele.Context) error {
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

	betsx, err := h.collapseBets(bets)
	if err != nil {
		return err
	}

	azartioBets := make([]*azartio.Bet, 0, len(betsx.Bets))
	for _, bet := range betsx.Bets {
		azartioBets = append(azartioBets, bet.ToAzartioBet(bet.UserID))
	}
	type Result struct {
		Results map[key]*azartio.RollResult
		Chats   map[int]*tele.Chat
		WonSign string
		Perks   map[int]pgtype.Varchar
	}
	considerPerks := func(result *azartio.RollResult, perks []string, amount int64) int64 {

		if result.Won {
			if h.d.Scope("+10%", perks...) {
				d := h.d.Get("mafioso")
				return amount + (amount/100)*int64(d.Int("multiplier"))
			}
		} else {
			if h.d.Scope("-10%", perks...) {
				d := h.d.Get("banker")
				return amount - (amount/100)*int64(d.Int("multiplier"))
			}
		}
		return amount
	}

	results, err := h.c.RollMany(azartioBets)
	if err != nil {
		return err
	}

	perks := make(map[int]pgtype.Varchar)
	jointResults := make(map[key]*azartio.RollResult)
	for _, result := range results {
		chat, err := h.b.ChatByID(strconv.Itoa(result.Bet.UserID))
		if _, ok := betsx.Chats[result.Bet.UserID]; !ok {
			if err != nil {
				return err
			}
			betsx.Chats[result.Bet.UserID] = chat

		}
		user, err := h.db.Users.ByID(chat)
		if err != nil {
			return err
		}
		perks[result.Bet.UserID] = user.Perks.Elements[rand.Intn(len(user.Perks.Elements)-1)]
		perksStrings := make([]string, len(user.Perks.Elements))
		for k, v := range user.Perks.Elements {
			perksStrings[k] = v.String
		}
		result.Amount = considerPerks(result, perksStrings, result.Amount)
		k := key{result.Bet.UserID, result.Bet.Sign}
		jointResults[k] = result

		if result.Won {
			err := h.db.Users.Charge(+result.Amount, betsx.Chats[result.Bet.UserID])
			if err != nil {
				return err
			}
		}
		err = h.db.Bets.MakeDoneByChat(result, betsx.Chats[result.Bet.UserID], c.Chat())
		if err != nil {
			return err
		}

	}

	err = h.db.Groups.UpdateState(c.Chat(), stateRolling)
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
		Chats:   betsx.Chats,
		Perks:   perks,
	}), h.lt.Markup(c, "roulette")

	time.AfterFunc(3*time.Second, func() {
		_, err := h.b.Edit(reservedMsg, a, b)
		if err != nil {
			h.b.OnError(err, c)
		}
		err = h.db.Groups.UpdateState(c.Chat(), stateIdle)
		if err != nil {
			h.b.OnError(err, c)
		}
	},
	)

	return err
}

func (h handler) OnPinned(c tele.Context) error {
	if c.Sender().ID == h.b.Me.ID {
		h.b.Delete(c.Message())
	}
	return nil
}

func (h handler) collapseBets(bets []*storage.Bet) (*Bets, error) {
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
