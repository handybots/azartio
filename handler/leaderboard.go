package handler

import (
	"strconv"

	"github.com/handybots/azartio/storage"
	tele "gopkg.in/tucnak/telebot.v3"
)

type Leaderboard struct {
	Top3  []LeaderboardUser
	Top10 []LeaderboardUser
	You   LeaderboardUser
}

type LeaderboardUser struct {
	storage.User
	*tele.Chat
	Place int
}

func (h handler) OnLeaderboard(c tele.Context) error {
	you, err := h.db.Users.ByID(c.Sender())
	if err != nil {
		return err
	}

	// NOTE: it will be extremely inefficient

	users, err := h.db.Users.Leaderboard()
	if err != nil {
		return err
	}

	var place int
	for i, user := range users {
		if user.ID == c.Sender().ID {
			place = i + 1
			break
		}
	}

	lb := Leaderboard{
		You: LeaderboardUser{User: you, Chat: c.Chat(), Place: place},
	}
	if len(users) > 10 {
		users = users[:9]
	}

	for i, user := range users {
		chat, err := h.b.ChatByID(strconv.Itoa(user.ID))
		if err != nil {
			return err
		}

		lbu := LeaderboardUser{
			User:  user,
			Chat:  chat,
			Place: i + 1,
		}

		if i < 3 {
			lb.Top3 = append(lb.Top3, lbu)
		} else {
			lb.Top10 = append(lb.Top10, lbu)
		}
	}

	return c.Send(h.lt.Text(c, "leaderboard", lb))
}
