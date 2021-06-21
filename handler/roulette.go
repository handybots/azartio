package handler

import (
	"github.com/handybots/azartio/azartio"
	"github.com/handybots/azartio/storage"
	"strconv"
	"time"

	tele "gopkg.in/tucnak/telebot.v3"
)

func (h handler) OnRoulette(c tele.Context) error {
	msg, err := h.b.Send(
		c.Chat(),
		h.lt.Text(c, "roulette"),
		h.lt.Markup(c, "roulette"),
	)
	if err != nil {
		return err
	}

	if !c.Message().Private() {
		h.b.Pin(msg, tele.Silent)
		c.Delete()
	}

	exists, err := h.db.Groups.Exists(c.Chat())
	if err != nil {
		return err
	}

	if exists {
		group, err := h.db.Groups.ByID(c.Chat())
		if err != nil {
			return err
		}

		if group.ID < 0 {
			h.b.Delete(tele.StoredMessage{
				MessageID: strconv.FormatInt(group.MessageID, 10),
				ChatID:    c.Chat().ID,
			})
		}

		// TODO: if the group has bets, edit the message to show them.
	} else if err := h.db.Groups.Create(c.Chat()); err != nil {
		return err
	}

	return h.db.Groups.UpdateMessage(c.Chat(), msg.ID)
}

func (h handler) OnRouletteGo(c tele.Context) error {
	defer c.Respond()

	group, err := h.db.Groups.ByID(c.Chat())
	if err != nil {
		return err
	}

	msg := tele.StoredMessage{
		MessageID: strconv.FormatInt(group.MessageID, 10),
		ChatID:    c.Chat().ID,
	}

	notDone, err := h.db.Bets.NotDoneByChat(c.Chat())
	if err != nil {
		return err
	}
	if len(notDone) == 0 {
		return c.Respond(&tele.CallbackResponse{Text: h.lt.Text(c, "no_bets")})
	}

	betsx, err := h.collapseBets(notDone)
	if err != nil {
		return err
	}

	bets := make([]azartio.Bet, 0, len(betsx.Bets))
	for _, bet := range betsx.Bets {
		bets = append(bets, azartio.Bet{
			Sign:   bet.Bet.Sign,
			Amount: bet.Bet.Amount,
			UserID: bet.UserID,
		})
	}

	results, err := h.rt.RollMany(bets)
	if err != nil {
		return err
	}

	joint := make(BetsMap)
	for i, result := range results {
		key := BetKey{result.Bet.UserID, result.Bet.Sign}
		joint[key] = BetOrResult{RollResult: results[i]}

		user, err := h.db.Users.ByID(c.Sender())
		if err != nil {
			return err
		}

		if result.Won {
			if h.dons.Scope("+10%", user.Perks()...) {
				result.Amount += int64(float64(result.Amount) * 0.1)
			}

			err = h.db.Users.Charge(result.Amount, betsx.Chats[result.Bet.UserID])
			if err != nil {
				return err
			}
		} else if h.dons.Scope("-10%", user.Perks()...) {
			cashback := int64(float64(result.Amount) * 0.1)

			err = h.db.Users.Charge(cashback, betsx.Chats[result.Bet.UserID])
			if err != nil {
				return err
			}
		}

		err = h.db.Bets.MakeDoneByChat(result, betsx.Chats[result.Bet.UserID], c.Chat())
		if err != nil {
			return err
		}
	}

	err = h.db.Groups.UpdateState(c.Chat(), storage.GroupStateRolling)
	if err != nil {
		return err
	}

	_, err = h.b.Edit(msg, h.lt.Text(c, "roulette_rolling", betsx))
	if err != nil {
		return err
	}

	time.Sleep(3 * time.Second)

	result := Bets{
		Sign:  results[0].Sign,
		Bets:  joint,
		Chats: betsx.Chats,
	}

	_, err = h.b.Edit(
		msg,
		h.lt.Text(c, "roulette_result", result),
		h.lt.Markup(c, "roulette"),
	)
	if err != nil {
		h.b.OnError(err, c)
	}

	return h.db.Groups.UpdateState(c.Chat(), storage.GroupStateNone)
}

func (h handler) OnPinned(c tele.Context) error {
	if c.Sender().ID == h.b.Me.ID {
		return h.b.Delete(c.Message())
	}
	return nil
}
