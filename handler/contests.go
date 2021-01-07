package handler

import (
	"strconv"

	"github.com/jackc/pgx/pgtype"

	"github.com/handybots/azartio/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

func (h *handler) OnMakeContest(c tele.Context) error {
	states[c.Sender().Recipient()] = StateAmountContest
	h.b.Send(c.Chat(), h.lt.Text(c, "enter_contest_amount"))
	return nil
}

func (h *handler) OnParticipate(c tele.Context) error {
	defer c.Respond()
	contest, err := h.db.Contests.NotDoneByUserID(c.Callback().Data)
	if err != nil {
		return err
	}
	if c.Sender().ID == contest.CreatorID {
		return nil
	}
	for _, v := range contest.Participants.Elements {
		if v.String == c.Sender().Recipient() {
			c.Respond(&tele.CallbackResponse{Text: h.lt.Text(c, "already_participating")})
			return nil
		}
	}
	err = h.db.Contests.AddParticipant(c.Callback().Data, c.Sender())
	if err != nil {
		c.Respond(&tele.CallbackResponse{Text: h.lt.Text(c, "cant_add_participant")})
		return err
	}
	contest.Participants.Elements = append(contest.Participants.Elements,
		pgtype.Varchar{c.Sender().Recipient(), '1'},
	)
	c.Respond(&tele.CallbackResponse{Text: h.lt.Text(c, "added_participant")})
	contestInfo, err := h.genContestInfo(contest)
	if err != nil {
		return err
	}
	_, err = h.b.Edit(c.Message(), h.lt.Text(c, "contest", contestInfo), h.lt.Markup(c, "contest", contest.CreatorID))
	return err
}

func (h *handler) OnContests(c tele.Context) error {
	return nil
}

func (h *handler) genContestInfo(contest storage.Contest) (Contest, error) {
	participants := make([]*tele.Chat, 0)
	for _, v := range contest.Participants.Elements {
		chat, err := h.b.ChatByID(v.String)
		if err != nil {
			continue
		}
		participants = append(participants, chat)
	}
	owner, err := h.b.ChatByID(strconv.Itoa(contest.CreatorID))
	if err != nil {
		return Contest{}, err
	}
	return Contest{owner, participants, contest.Amount}, nil
}
