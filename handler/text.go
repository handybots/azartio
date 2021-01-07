package handler

import (
	"errors"
	"regexp"
	"strconv"

	tele "gopkg.in/tucnak/telebot.v3"
)

var reBet = regexp.MustCompile(`(\d+)\s*(?:на\s*)?(к(?:расное)?|ч(?:[её]рное)?|з(?:ел[её]ное)?)$`)

const (
	StateAmountContest = "amount_contest"
)

var states = make(map[string]string)

type Contest struct {
	Creator      *tele.Chat
	Participants []*tele.Chat
	Amount       int64
}

func (h handler) OnText(c tele.Context) error {
	if c.Message().Text == h.lt.Text(c, "my_balance") {
		return h.OnBalance(c)
	}
	if reBet.MatchString(c.Message().Text) {
		return h.OnBet(c)
	}

	state, ok := states[c.Sender().Recipient()]
	if !ok {
		return errors.New("handler: state not found")
	}

	switch {
	case state == StateAmountContest:
		amount, err := strconv.ParseInt(c.Message().Text, 10, 64)
		if err != nil {
			return err
		}
		delete(states, c.Sender().Recipient())

		err = h.db.Contests.Create(c.Sender(), c.Chat(), amount)
		if err != nil {
			c.Send(h.lt.Text(c, "something_went_wrong"))
			return err
		}

		contest, err := h.db.Contests.NotDoneByUserID(c.Sender().Recipient())
		if err != nil {
			return err
		}

		contestInfo, err := h.genContestInfo(contest)
		if err != nil {
			return err
		}
		owner, err := h.b.ChatByID(strconv.Itoa(contest.CreatorID))
		if err != nil {
			return err
		}
		h.b.Send(c.Chat(), h.lt.Text(c, "contest_created"))
		h.b.Send(c.Chat(),
			h.lt.Text(c, "contest", contestInfo),
			h.lt.Markup(c, "contest", owner.Recipient()))
	}
	return nil
}
