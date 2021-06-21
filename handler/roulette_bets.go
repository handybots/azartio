package handler

import (
	"database/sql"
	"errors"
	"github.com/handybots/azartio/azartio"
	"github.com/handybots/azartio/storage"
	tele "gopkg.in/tucnak/telebot.v3"
	"strconv"
)

// TODO: refactor

type (
	BetKey struct {
		UserID int
		Sign   string
	}

	BetOrResult struct {
		*storage.Bet
		*azartio.RollResult
	}

	BetsMap = map[BetKey]BetOrResult

	Bets struct {
		Sign  string
		Bets  BetsMap
		Chats map[int]*tele.Chat
	}
)

func (h handler) OnRouletteBet(c tele.Context) error {
	defer c.Respond()

	group, err := h.db.Groups.ByID(c.Chat())
	if err != nil {
		return err
	}

	msg := tele.StoredMessage{
		MessageID: strconv.Itoa(int(group.MessageID)),
		ChatID:    c.Chat().ID,
	}

	var amount int64
	sign := c.Callback().Unique

	if c.Data() == "x2" {
		notDone, err := h.db.Bets.NotDoneByUserID(c.Sender(), c.Chat())
		if err == sql.ErrNoRows {
			return c.Respond(&tele.CallbackResponse{
				Text: h.lt.Text(c, "did_not_bet"),
			})
		} else if err != nil {
			return err
		}

		betsx, err := h.collapseBets(notDone)
		if err != nil {
			return err
		}

		for _, bet := range betsx.Bets {
			if bet.Bet.Sign == sign {
				amount = bet.Bet.Amount
				break
			}
		}
	} else {
		amount, _ = strconv.ParseInt(c.Data(), 10, 64)
	}

	if amount == 0 {
		return c.Respond(&tele.CallbackResponse{
			Text: h.lt.Text(c, "did_not_bet"),
		})
	}

	balance, err := h.db.Users.Balance(c.Sender())
	if err != nil {
		return err
	}

	if balance-amount < 0 {
		return c.Respond(&tele.CallbackResponse{
			Text: h.lt.Text(c, "not_enough_money"),
		})
	}

	if _, ok := azartio.Colors[sign]; !ok {
		return errors.New("handler/bet: unknown bet sign")
	}

	if err := h.db.Users.Charge(-amount, c.Sender()); err != nil {
		return err
	}

	bet := azartio.NewBet(sign, amount, c.Sender().ID)
	if err = h.db.Bets.Create(bet, c.Chat()); err != nil {
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
		msg,
		h.lt.Text(c, "roulette_bets", betsx),
		h.lt.Markup(c, "roulette"),
	)
	return err
}

func (h handler) collapseBets(bets []storage.Bet) (*Bets, error) {
	chats := make(map[int]*tele.Chat)
	joint := make(BetsMap)

	for i, bet := range bets {
		if _, ok := chats[bet.UserID]; !ok {
			chat, err := h.b.ChatByID(strconv.Itoa(bet.UserID))
			if err != nil {
				return nil, err
			}
			chats[bet.UserID] = chat
		}

		k := BetKey{bet.UserID, bet.Sign}
		if b, ok := joint[k]; ok {
			b.Bet.Amount += bet.Amount
		} else {
			joint[k] = BetOrResult{Bet: &bets[i]}
		}
	}

	return &Bets{Bets: joint, Chats: chats}, nil
}

func (bets Bets) Sorted() []BetsMap {
	return []BetsMap{
		bets.BySign(azartio.Red),
		bets.BySign(azartio.Black),
		bets.BySign(azartio.Clever),
	}
}

func (bets Bets) BySign(sign string) BetsMap {
	sorted := make(BetsMap)
	for key, bet := range bets.Bets {
		if key.Sign == sign {
			sorted[key] = bet
		}
	}
	return sorted
}
